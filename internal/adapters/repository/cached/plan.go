package cached

import (
	"context"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type CachedPlanRepo struct {
	repo  ports.PlanRepository
	cache ports.Cache
	ttl   time.Duration
}

func NewCachedPlanRepo(repo ports.PlanRepository, cache ports.Cache, ttl time.Duration) ports.PlanRepository {
	return &CachedPlanRepo{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (r *CachedPlanRepo) FindByID(_ context.Context, id string) (*domain.Plan, error) {

}

func (r *CachedPlanRepo) FindDefault(ctx context.Context) (*domain.Plan, error) {
	return r.FindByID(ctx, "free")
}

func (r *CachedPlanRepo) FindAll(_ context.Context) ([]*domain.Plan, error) {

}

func (r *CachedPlanRepo) Save(_ context.Context, plan *domain.Plan) error {

}
