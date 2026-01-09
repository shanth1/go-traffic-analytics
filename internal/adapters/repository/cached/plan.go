package cached

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type PlanRepo struct {
	repo  ports.PlanRepository
	cache ports.Cache
	ttl   time.Duration
}

func NewPlanRepo(repo ports.PlanRepository, cache ports.Cache, ttl time.Duration) ports.PlanRepository {
	return &PlanRepo{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (r *PlanRepo) buildKey(id string) string {
	return fmt.Sprintf("gotrace:plan:id:%s", id)
}

func (r *PlanRepo) buildAllKey() string {
	return "gotrace:plan:all"
}

func (r *PlanRepo) FindByID(ctx context.Context, id string) (*domain.Plan, error) {
	const op = "cached.PlanRepo.FindByID"

	key := r.buildKey(id)

	val, err := r.cache.Get(ctx, key)
	if err == nil {
		if bytesVal, ok := val.([]byte); ok {
			var plan domain.Plan
			if jsonErr := json.Unmarshal(bytesVal, &plan); jsonErr == nil {
				return &plan, nil
			}
			// TODO: logging (unmarshal error)
		}
	}
	// else if !errors.Is(err, ports.ErrCacheMiss) {
	// TODO: logging (cache error)
	// }

	plan, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ops.E(op, err)
	}

	go func() {
		bytes, err := json.Marshal(plan)
		if err != nil {
			// TODO: logging (marshal error)
			return
		}

		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// TODO: logging (failed to set cache)
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return plan, nil
}

func (r *PlanRepo) FindDefault(ctx context.Context) (*domain.Plan, error) {
	const op = "cached.PlanRepo.FindDefault"

	plan, err := r.FindByID(ctx, "free")
	if err != nil {
		return nil, ops.E(op, err)
	}

	return plan, nil
}

func (r *PlanRepo) FindAll(ctx context.Context) ([]*domain.Plan, error) {
	const op = "cached.PlanRepo.FindAll"

	key := r.buildAllKey()

	val, err := r.cache.Get(ctx, key)
	if err == nil {
		if bytesVal, ok := val.([]byte); ok {
			var plans []*domain.Plan
			if jsonErr := json.Unmarshal(bytesVal, &plans); jsonErr == nil {
				return plans, nil
			}
			// TODO: logging (unmarshal error)
		}
	}
	// else if !errors.Is(err, ports.ErrCacheMiss) {
	// TODO: logging (cache error)
	// }

	plans, err := r.repo.FindAll(ctx)
	if err != nil {
		return nil, ops.E(op, err)
	}

	go func() {
		bytes, err := json.Marshal(plans)
		if err != nil {
			// TODO: logging (marshal error)
			return
		}

		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// TODO: logging (failed to set cache)
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return plans, nil
}

func (r *PlanRepo) Save(ctx context.Context, plan *domain.Plan) error {
	const op = "cached.PlanRepo.Save"

	if err := r.repo.Save(ctx, plan); err != nil {
		return ops.E(op, err)
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// TODO: logging
		_ = r.cache.Delete(bgCtx, r.buildKey(plan.ID))
		_ = r.cache.Delete(bgCtx, r.buildAllKey())
	}()

	return nil
}
