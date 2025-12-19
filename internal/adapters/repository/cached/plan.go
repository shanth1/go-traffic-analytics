package cached

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (r *CachedPlanRepo) buildKey(id string) string {
	return fmt.Sprintf("gotrace:plan:id:%s", id)
}

func (r *CachedPlanRepo) buildAllKey() string {
	return "gotrace:plan:all"
}

func (r *CachedPlanRepo) FindByID(ctx context.Context, id string) (*domain.Plan, error) {
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
	} else if !errors.Is(err, ports.ErrCacheMiss) {
		// TODO: logging (cache error)
	}

	plan, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		bytes, err := json.Marshal(plan)
		if err != nil {
			// TODO: logging (marshal error)
			return
		}

		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.cache.Set(bgCtx, key, bytes, r.ttl); err != nil {
			// TODO: logging (failed to set cache)
		}
	}()

	return plan, nil
}

func (r *CachedPlanRepo) FindDefault(ctx context.Context) (*domain.Plan, error) {
	return r.FindByID(ctx, "free")
}

func (r *CachedPlanRepo) FindAll(ctx context.Context) ([]*domain.Plan, error) {
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
	} else if !errors.Is(err, ports.ErrCacheMiss) {
		// TODO: logging (cache error)
	}

	plans, err := r.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		bytes, err := json.Marshal(plans)
		if err != nil {
			// TODO: logging (marshal error)
			return
		}

		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.cache.Set(bgCtx, key, bytes, r.ttl); err != nil {
			// TODO: logging (failed to set cache)
		}
	}()

	return plans, nil
}

func (r *CachedPlanRepo) Save(ctx context.Context, plan *domain.Plan) error {
	if err := r.repo.Save(ctx, plan); err != nil {
		return err
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.cache.Delete(bgCtx, r.buildKey(plan.ID)); err != nil {
			// TODO: logging
		}

		if err := r.cache.Delete(bgCtx, r.buildAllKey()); err != nil {
			// TODO: logging
		}
	}()

	return nil
}
