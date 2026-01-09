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

type CampaignRepo struct {
	repo  ports.CampaignRepository
	cache ports.Cache
	ttl   time.Duration
}

func NewCampaignRepo(repo ports.CampaignRepository, cache ports.Cache, ttl time.Duration) ports.CampaignRepository {
	return &CampaignRepo{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (r *CampaignRepo) buildKey(id string) string {
	return fmt.Sprintf("gotrace:campaign:id:%s", id)
}

func (r *CampaignRepo) Save(ctx context.Context, camp *domain.Campaign) error {
	const op = "cached.CampaignRepo.Save"

	if err := r.repo.Save(ctx, camp); err != nil {
		return ops.E(op, err)
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// TODO: logging
		_ = r.cache.Delete(bgCtx, r.buildKey(camp.ID))
	}()

	return nil
}

func (r *CampaignRepo) FindByID(ctx context.Context, id string) (*domain.Campaign, error) {
	const op = "cached.CampaignRepo.FindByID"

	key := r.buildKey(id)
	val, err := r.cache.Get(ctx, key)
	if err == nil {
		if bytesVal, ok := val.([]byte); ok {
			var camp domain.Campaign
			if jsonErr := json.Unmarshal(bytesVal, &camp); jsonErr == nil {
				return &camp, nil
			}
			// TODO: logging (unmarshal error)
		}
	}
	// else if !errors.Is(err, ports.ErrCacheMiss) {
	// TODO: logging (cache error)
	// }

	camp, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ops.E(op, err)
	}

	go func() {
		bytes, err := json.Marshal(camp)
		if err == nil {
			bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			// TODO: logging
			_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
		}
		// TODO: logging
	}()

	return camp, nil
}

func (r *CampaignRepo) FindAll(ctx context.Context, filter domain.CampaignFilter) ([]*domain.Campaign, error) {
	const op = "cached.CampaignRepo.FindAll"

	campaigns, err := r.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, ops.E(op, err)
	}

	return campaigns, nil
}

func (r *CampaignRepo) Count(ctx context.Context, filter domain.CampaignFilter) (int64, error) {
	const op = "cached.CampaignRepo.Count"

	count, err := r.repo.Count(ctx, filter)
	if err != nil {
		return 0, ops.E(op, err)
	}

	return count, nil
}

func (r *CampaignRepo) Delete(ctx context.Context, userID domain.UserID, id string) error {
	const op = "cached.CampaignRepo.Delete"

	if err := r.repo.Delete(ctx, userID, id); err != nil {
		return ops.E(op, err)
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// TODO: logging
		_ = r.cache.Delete(bgCtx, r.buildKey(id))
	}()

	return nil
}
