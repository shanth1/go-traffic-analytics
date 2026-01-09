package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"golang.org/x/sync/errgroup"
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
		count, err := s.linkRepo.Count(ctx, domain.LinkFilter{UserID: cmd.UserID})
		if err != nil {
			return nil, err
		}
		if int(count) >= plan.MaxLinks {
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

func (s *LinkService) GetLinkList(ctx context.Context, filter domain.LinkFilter) ([]*domain.Link, int64, error) {
	var (
		links []*domain.Link
		total int64
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		links, err = s.linkRepo.FindAll(gCtx, filter)
		return err
	})

	g.Go(func() error {
		var err error
		total, err = s.linkRepo.Count(gCtx, filter)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	return links, total, nil
}

// TODO:
// soft delete? isActive = true
func (s *LinkService) DeleteLink(ctx context.Context, userID domain.UserID, id domain.LinkID) error {
	return s.linkRepo.Delete(ctx, userID, id)
}
