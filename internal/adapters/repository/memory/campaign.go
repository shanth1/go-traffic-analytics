package memoryrepo

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/shanth1/gotools/errs"
	"github.com/shanth1/gotools/ops"
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

func (r *InMemoryCampaignRepo) Save(_ context.Context, camp *domain.Campaign) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.campaigns[camp.ID] = camp
	return nil
}

func (r *InMemoryCampaignRepo) FindByID(_ context.Context, id string) (*domain.Campaign, error) {
	const op = "memoryrepo.InMemoryCampaignRepo.FindByID"

	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.campaigns[id]
	if !ok {
		return nil, ops.E(op, ops.KindNotFound, fmt.Errorf("campaign with id %q: %w", id, errs.ErrNotFound))
	}

	return c, nil
}

func (r *InMemoryCampaignRepo) FindAll(_ context.Context, filter domain.CampaignFilter) ([]*domain.Campaign, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var matches []*domain.Campaign
	for _, c := range r.campaigns {
		if c.UserID == filter.UserID {
			matches = append(matches, c)
		}
	}

	// Sort by creation desc
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].CreatedAt.After(matches[j].CreatedAt)
	})

	// Pagination
	if filter.Offset >= len(matches) {
		return []*domain.Campaign{}, nil
	}

	res := matches[filter.Offset:]
	if filter.Limit > 0 && filter.Limit < len(res) {
		res = res[:filter.Limit]
	}

	return res, nil
}

func (r *InMemoryCampaignRepo) Count(_ context.Context, filter domain.CampaignFilter) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var count int64
	for _, c := range r.campaigns {
		if c.UserID == filter.UserID {
			count++
		}
	}
	return count, nil
}

func (r *InMemoryCampaignRepo) Delete(_ context.Context, userID domain.UserID, id string) error {
	const op = "memoryrepo.InMemoryCampaignRepo.Delete"

	r.mu.Lock()
	defer r.mu.Unlock()

	c, ok := r.campaigns[id]
	if !ok {
		return ops.E(op, ops.KindNotFound, fmt.Errorf("campaign with id %q: %w", id, errs.ErrNotFound))
	}

	if c.UserID != userID {
		return ops.E(op, ops.KindUnauthorized, fmt.Errorf("expected user with id %q, got %q: %w", userID, c.UserID, errs.ErrUnauthorized))
	}

	delete(r.campaigns, id)

	return nil
}
