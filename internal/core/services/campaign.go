package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
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

func (s *CampaignService) GetCampaigns(ctx context.Context, userID domain.UserID) ([]*domain.Campaign, error) {
	return s.repo.FindAllByUserID(ctx, userID)
}

func (s *CampaignService) DeleteCampaign(ctx context.Context, userID domain.UserID, id string) error {
	return s.repo.Delete(ctx, userID, id)
}
