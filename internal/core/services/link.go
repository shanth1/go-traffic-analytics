package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

var ErrLimitReached = errors.New("plan limit reached: upgrade your plan to create more links")

type LinkService struct {
	linkRepo ports.LinkRepository
	userRepo ports.UserRepository
	planRepo ports.PlanRepository
}

func NewLinkService(lr ports.LinkRepository, ur ports.UserRepository, pr ports.PlanRepository) *LinkService {
	return &LinkService{
		linkRepo: lr,
		userRepo: ur,
		planRepo: pr,
	}
}

// --- Links ---

func (s *LinkService) CreateLink(ctx context.Context, cmd domain.CreateLinkCmd) (*domain.Link, error) {
	user, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	plan, err := s.planRepo.FindByID(ctx, user.PlanID)
	if err != nil {
		return nil, err
	}

	if plan.MaxLinks != -1 {
		currentCount, err := s.linkRepo.CountByUserID(ctx, cmd.UserID)
		if err != nil {
			return nil, err
		}
		if int(currentCount) >= plan.MaxLinks {
			return nil, ErrLimitReached
		}
	}

	// Slug Generation
	finalSlug := cmd.CustomSlug
	if finalSlug == "" {
		finalSlug = uuid.New().String()[:8] // Simple random slug
	}

	link := &domain.Link{
		ID:         domain.LinkID(uuid.NewString()),
		UserID:     cmd.UserID,
		CampaignID: cmd.CampaignID,
		TargetURL:  cmd.TargetURL,
		Slug:       finalSlug,
		IsActive:   true,
		CreatedAt:  time.Now(),
	}

	if err := s.linkRepo.Save(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

func (s *LinkService) GetLinkList(ctx context.Context, filter domain.LinkFilter) ([]*domain.Link, error) {
	return s.linkRepo.FindAll(ctx, filter)
}

// TODO:
// soft delete? isActive = true
func (s *LinkService) DeleteLink(_ context.Context, _ domain.LinkID) error {
	return nil
}
