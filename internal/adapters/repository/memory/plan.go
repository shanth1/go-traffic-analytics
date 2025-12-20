package memoryrepo

import (
	"context"
	"errors"
	"sync"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type InMemoryPlanRepo struct {
	mu    sync.RWMutex
	plans map[string]*domain.Plan
}

func NewPlanRepo() ports.PlanRepository {
	repo := &InMemoryPlanRepo{
		plans: make(map[string]*domain.Plan),
	}

	return repo
}

func (r *InMemoryPlanRepo) FindByID(_ context.Context, id string) (*domain.Plan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.plans[id]
	if !ok {
		return nil, errors.New("plan not found")
	}
	return p, nil
}

func (r *InMemoryPlanRepo) FindDefault(ctx context.Context) (*domain.Plan, error) {
	return r.FindByID(ctx, "free")
}

func (r *InMemoryPlanRepo) FindAll(_ context.Context) ([]*domain.Plan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Plan, 0, len(r.plans))
	for _, p := range r.plans {
		res = append(res, p)
	}
	return res, nil
}

func (r *InMemoryPlanRepo) Save(_ context.Context, plan *domain.Plan) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plans[plan.ID] = plan
	return nil
}
