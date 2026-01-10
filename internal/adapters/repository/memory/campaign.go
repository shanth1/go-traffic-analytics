package memory

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

type CampaignRepo struct {
	mu        sync.RWMutex
	campaigns map[string]*domain.Campaign
}

func NewCampaignRepo() ports.CampaignRepository {
	return &CampaignRepo{
		campaigns: make(map[string]*domain.Campaign),
	}
}

func (r *CampaignRepo) Save(_ context.Context, camp *domain.Campaign) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.campaigns[camp.ID] = camp
	return nil
}

func (r *CampaignRepo) FindByID(_ context.Context, id string) (*domain.Campaign, error) {
	const op = "memory.CampaignRepo.FindByID"

	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.campaigns[id]
	if !ok {
		return nil, ops.WrapMsg(op, ops.KindNotFound, errs.ErrNotFound, fmt.Sprintf("%q campaign not found", id))
	}

	return c, nil
}

func (r *CampaignRepo) FindAll(_ context.Context, filter domain.CampaignFilter) ([]*domain.Campaign, error) {
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

func (r *CampaignRepo) Count(_ context.Context, filter domain.CampaignFilter) (int64, error) {
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

func (r *CampaignRepo) Delete(_ context.Context, userID domain.UserID, id string) error {
	const op = "memory.CampaignRepo.Delete"

	r.mu.Lock()
	defer r.mu.Unlock()

	c, ok := r.campaigns[id]
	if !ok {
		return ops.WrapMsg(op, ops.KindNotFound, errs.ErrNotFound, fmt.Sprintf("%q campaign not found", id))
	}

	if c.UserID != userID {
		techErr := fmt.Errorf("user %q attempted to access campaign %q owned by %q: %w",
			userID, id, c.UserID, errs.ErrForbidden)
		return ops.WrapMsg(op, ops.KindPermission, techErr, "access denied")
	}

	delete(r.campaigns, id)

	return nil
}
