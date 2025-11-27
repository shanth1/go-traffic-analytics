package memory

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
	// Инициализация базовых тарифов при создании репо
	repo.bootstrapPlans()
	return repo
}

func (r *InMemoryPlanRepo) bootstrapPlans() {
	// 1. Free
	r.plans["free"] = &domain.Plan{
		ID:             "free",
		Name:           "Starter Free",
		PriceCents:     0,
		MaxLinks:       10,
		MaxClicksMonth: 1000,
		CanExportData:  false,
		IsActive:       true,
	}
	// 2. Pro
	r.plans["pro"] = &domain.Plan{
		ID:             "pro",
		Name:           "Professional",
		PriceCents:     1900, // $19.00
		MaxLinks:       100,
		MaxClicksMonth: 50000,
		CanExportData:  true,
		IsActive:       true,
	}
	// 3. Enterprise
	r.plans["enterprise"] = &domain.Plan{
		ID:             "enterprise",
		Name:           "Business Unlimited",
		PriceCents:     9900,
		MaxLinks:       -1, // Безлимит
		MaxClicksMonth: 1000000,
		CanExportData:  true,
		IsActive:       true,
	}
}

func (r *InMemoryPlanRepo) FindByID(ctx context.Context, id string) (*domain.Plan, error) {
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

func (r *InMemoryPlanRepo) FindAll(ctx context.Context) ([]*domain.Plan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []*domain.Plan
	for _, p := range r.plans {
		res = append(res, p)
	}
	return res, nil
}

func (r *InMemoryPlanRepo) Save(ctx context.Context, plan *domain.Plan) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plans[plan.ID] = plan
	return nil
}
