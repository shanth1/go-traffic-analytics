package cached

import (
	"context"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type CachedCampaignRepo struct {
	repo  ports.CampaignRepository
	cache ports.Cache
	ttl   time.Duration
}

func NewCampaignRepo(repo ports.CampaignRepository, cache ports.Cache, ttl time.Duration) ports.CampaignRepository {
	return &CachedCampaignRepo{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (r *CachedCampaignRepo) Save(_ context.Context, camp *domain.Campaign) error {

}

func (r *CachedCampaignRepo) FindByID(_ context.Context, id string) (*domain.Campaign, error) {

}

func (r *CachedCampaignRepo) FindAllByUserID(_ context.Context, userID string) ([]*domain.Campaign, error) {

}

func (r *CachedCampaignRepo) Delete(_ context.Context, id string) error {

}
