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
	campRepo ports.CampaignRepository
	userRepo ports.UserRepository
	planRepo ports.PlanRepository
}

func NewLinkService(l ports.LinkRepository, c ports.CampaignRepository, u ports.UserRepository, p ports.PlanRepository) *LinkService {
	return &LinkService{linkRepo: l, campRepo: c, userRepo: u, planRepo: p}
}

// --- Campaigns ---

func (s *LinkService) GetUserCampaigns(ctx context.Context, userID string) ([]*domain.Campaign, error) {
	return s.campRepo.FindAllByUserID(ctx, userID)
}

func (s *LinkService) CreateCampaign(ctx context.Context, userID, name string) (*domain.Campaign, error) {
	camp := &domain.Campaign{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
	}
	if err := s.campRepo.Save(ctx, camp); err != nil {
		return nil, err
	}
	return camp, nil
}

// --- Links ---

func (s *LinkService) GetLinks(ctx context.Context, campaignID string) ([]*domain.Link, error) {
	return s.linkRepo.FindAllByCampaignID(ctx, campaignID)
}

func (s *LinkService) CreateLink(ctx context.Context, userID, campaignID, targetURL, customSlug string) (*domain.Link, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	plan, err := s.planRepo.FindByID(ctx, user.PlanID)
	if err != nil {
		return nil, err
	}

	if plan.MaxLinks != -1 {
		currentCount, err := s.linkRepo.CountByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if int(currentCount) >= plan.MaxLinks {
			return nil, ErrLimitReached
		}
	}

	// Slug Generation
	finalSlug := customSlug
	if finalSlug == "" {
		finalSlug = uuid.New().String()[:8] // Simple random slug
	}

	link := &domain.Link{
		ID:         uuid.New().String(),
		UserID:     userID,
		CampaignID: campaignID,
		TargetURL:  targetURL,
		Slug:       finalSlug,
		IsActive:   true,
		CreatedAt:  time.Now(),
	}

	if err := s.linkRepo.Save(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

// TODO:
// soft delete? isActive = true
func (s *LinkService) DeleteLink(ctx context.Context, id string) error {
	return nil
}
