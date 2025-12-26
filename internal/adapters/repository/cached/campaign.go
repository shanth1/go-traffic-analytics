package cached

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	if err := r.repo.Save(ctx, camp); err != nil {
		return err
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		key := r.buildKey(camp.ID)

		// TODO: logging (failed to clear cache)
		_ = r.cache.Delete(bgCtx, key)
	}()

	return nil
}

func (r *CampaignRepo) FindByID(ctx context.Context, id string) (*domain.Campaign, error) {
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
		return nil, err
	}

	go func() {
		bytes, err := json.Marshal(camp)
		if err != nil {
			// TODO: logging (marshal error)
			return
		}

		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// TODO: logging (failed to set cache)
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return camp, nil
}

func (r *CampaignRepo) FindAllByUserID(ctx context.Context, userID domain.UserID) ([]*domain.Campaign, error) {
	return r.repo.FindAllByUserID(ctx, userID)
}

func (r *CampaignRepo) Delete(ctx context.Context, id string) error {
	if err := r.repo.Delete(ctx, id); err != nil {
		return err
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		key := r.buildKey(id)

		// TODO: logging (failed to delete from cache)
		_ = r.cache.Delete(bgCtx, key)
	}()

	return nil
}
