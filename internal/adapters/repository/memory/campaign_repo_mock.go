package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type InMemoryCampaignRepo struct {
	mu        sync.RWMutex
	campaigns map[string]*domain.Campaign
}

func NewCampaignRepo() ports.CampaignRepository {
	return &InMemoryCampaignRepo{
		campaigns: make(map[string]*domain.Campaign),
	}
}

func (r *InMemoryCampaignRepo) Save(ctx context.Context, camp *domain.Campaign) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.campaigns[camp.ID] = camp
	return nil
}

func (r *InMemoryCampaignRepo) FindByID(ctx context.Context, id string) (*domain.Campaign, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.campaigns[id]
	if !ok {
		return nil, errors.New("campaign not found")
	}
	return c, nil
}

func (r *InMemoryCampaignRepo) FindAllByUserID(ctx context.Context, userID string) ([]*domain.Campaign, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.Campaign
	for _, c := range r.campaigns {
		if c.UserID == userID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *InMemoryCampaignRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.campaigns, id)
	return nil
}
