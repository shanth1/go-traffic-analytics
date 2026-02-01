package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"golang.org/x/sync/errgroup"
)

type LinkService struct {
	linkRepo ports.LinkRepository
	userRepo ports.UserRepository
	planRepo ports.PlanRepository
	logger   log.Logger
}

func NewLinkService(
	lr ports.LinkRepository,
	ur ports.UserRepository,
	pr ports.PlanRepository,
	logger log.Logger,
) *LinkService {
	return &LinkService{
		linkRepo: lr,
		userRepo: ur,
		planRepo: pr,
		logger:   logger,
	}
}

func (s *LinkService) CreateLink(ctx context.Context, cmd domain.CreateLinkCmd) (*domain.Link, error) {
	const op = "LinkService.CreateLink"

	user, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, ops.WrapMsg(op, ops.KindInternal, err, "failed to fetch user")
	}

	plan, err := s.planRepo.FindByID(ctx, user.PlanID)
	if err != nil {
		return nil, ops.WrapMsg(op, ops.KindInternal, err, "failed to fetch plan")
	}

	if plan.MaxLinks != -1 {
		count, err := s.linkRepo.Count(ctx, domain.LinkFilter{UserID: cmd.UserID})
		if err != nil {
			return nil, ops.WrapMsg(op, ops.KindInternal, err, "failed to count user links")
		}
		if int(count) >= plan.MaxLinks {
			return nil, ops.WrapMsg(op, ops.KindPermission, errors.New("plan limit reached"), "upgrade your plan to create more links")
		}
	}

	finalSlug := cmd.CustomSlug
	if finalSlug == "" {
		// TODO: Retry mechanism on collision
		finalSlug = uuid.New().String()[:7]
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
		return nil, ops.WrapMsg(op, ops.KindInternal, err, "failed to save link")
	}

	return link, nil
}

func (s *LinkService) GetLinkList(ctx context.Context, filter domain.LinkFilter) ([]*domain.Link, int64, error) {
	const op = "LinkService.GetLinkList"

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
		return nil, 0, ops.WrapMsg(op, ops.KindInternal, err, "failed to fetch links")
	}

	return links, total, nil
}

func (s *LinkService) DeleteLink(ctx context.Context, userID domain.UserID, id domain.LinkID) error {
	const op = "LinkService.DeleteLink"
	if err := s.linkRepo.Delete(ctx, userID, id); err != nil {
		return ops.WrapMsg(op, ops.KindInternal, err, "failed to delete link")
	}
	return nil
}
