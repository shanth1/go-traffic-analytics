package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"golang.org/x/sync/errgroup"
)

type CampaignService struct {
	repo     ports.CampaignRepository
	linkRepo ports.LinkRepository
	tx       ports.Transactor
	logger   log.Logger
}

func NewCampaignService(
	campaignRepo ports.CampaignRepository,
	linkRepo ports.LinkRepository,
	tx ports.Transactor,
	logger log.Logger,
) *CampaignService {
	return &CampaignService{
		repo:     campaignRepo,
		linkRepo: linkRepo,
		tx:       tx,
		logger:   logger,
	}
}

func (s *CampaignService) CreateCampaign(ctx context.Context, userID domain.UserID, name string) (*domain.Campaign, error) {
	const op = "CampaignService.CreateCampaign"

	camp := &domain.Campaign{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Save(ctx, camp); err != nil {
		return nil, ops.WrapMsg(op, ops.KindInternal, err, "failed to save campaign")
	}

	return camp, nil
}

func (s *CampaignService) GetCampaigns(ctx context.Context, filter domain.CampaignFilter) ([]*domain.Campaign, int64, error) {
	const op = "CampaignService.GetCampaigns"

	var (
		campaigns []*domain.Campaign
		total     int64
	)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		campaigns, err = s.repo.FindAll(gCtx, filter)
		return err
	})

	g.Go(func() error {
		var err error
		total, err = s.repo.Count(gCtx, filter)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, 0, ops.WrapMsg(op, ops.KindInternal, err, "failed to fetch campaigns")
	}

	return campaigns, total, nil
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, userID domain.UserID, id string) error {
	const op = "CampaignService.DeleteCampaign"

	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.Delete(txCtx, userID, id); err != nil {
			return err
		}

		if err := s.linkRepo.DeleteByCampaignID(txCtx, id); err != nil {
			return fmt.Errorf("failed to cascade delete links: %w", err)
		}

		return nil
	})

	if err != nil {
		return ops.WrapMsg(op, ops.KindInternal, err, "failed to delete campaign")
	}

	return nil
}
