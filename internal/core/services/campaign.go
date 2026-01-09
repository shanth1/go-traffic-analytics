package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"golang.org/x/sync/errgroup"
)

type CampaignService struct {
	repo ports.CampaignRepository
}

func NewCampaignService(campaignRepo ports.CampaignRepository) *CampaignService {
	return &CampaignService{
		repo: campaignRepo,
	}
}

func (s *CampaignService) CreateCampaign(ctx context.Context, userID domain.UserID, name string) (*domain.Campaign, error) {
	camp := &domain.Campaign{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Save(ctx, camp); err != nil {
		return nil, err
	}
	return camp, nil
}

func (s *CampaignService) GetCampaigns(ctx context.Context, filter domain.CampaignFilter) ([]*domain.Campaign, int64, error) {
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
		return nil, 0, err
	}

	return campaigns, total, nil
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, userID domain.UserID, id string) error {
	return s.repo.Delete(ctx, userID, id)
}
