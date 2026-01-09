package memoryrepo

import (
	"context"
	"fmt"
	"sync"

	"github.com/shanth1/gotools/errs"
	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type MemoryPlanRepo struct {
	mu    sync.RWMutex
	plans map[string]*domain.Plan
}

func NewPlanRepo() ports.PlanRepository {
	repo := &MemoryPlanRepo{
		plans: make(map[string]*domain.Plan),
	}

	return repo
}

func (r *MemoryPlanRepo) FindByID(_ context.Context, id string) (*domain.Plan, error) {
	const op = "memoryrepo.MemoryPlanRepo.FindByID"

	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.plans[id]
	if !ok {
		return nil, ops.E(op, ops.KindNotFound, fmt.Errorf("plan with id %q: %w", id, errs.ErrNotFound))
	}

	return p, nil
}

func (r *MemoryPlanRepo) FindDefault(ctx context.Context) (*domain.Plan, error) {
	return r.FindByID(ctx, "free")
}

func (r *MemoryPlanRepo) FindAll(_ context.Context) ([]*domain.Plan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*domain.Plan, 0, len(r.plans))
	for _, p := range r.plans {
		res = append(res, p)
	}
	return res, nil
}

func (r *MemoryPlanRepo) Save(_ context.Context, plan *domain.Plan) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plans[plan.ID] = plan
	return nil
}
