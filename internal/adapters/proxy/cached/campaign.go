package cachedproxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotools/logkeys"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type CampaignRepo struct {
	repo   ports.CampaignRepository
	cache  ports.Cache
	logger log.Logger
	ttl    time.Duration
}

func NewCampaignRepo(repo ports.CampaignRepository, cache ports.Cache, logger log.Logger, ttl time.Duration) ports.CampaignRepository {
	return &CampaignRepo{
		repo:   repo,
		cache:  cache,
		logger: logger,
		ttl:    ttl,
	}
}

func (r *CampaignRepo) buildKey(id string) string {
	return fmt.Sprintf("gotrace:campaign:id:%s", id)
}

func (r *CampaignRepo) Save(ctx context.Context, camp *domain.Campaign) error {
	if err := r.repo.Save(ctx, camp); err != nil {
		return err
	}

	r.invalidate(camp.ID)

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
			r.logger.Error().Err(err).Str(logkeys.CacheKey, key).Msg("failed to unmarshal campaign from cache")
		}
	} else if !errors.Is(err, ports.ErrCacheMiss) {
		r.logger.Error().Err(err).Str(logkeys.CacheKey, key).Msg("cache miss or error")
	}

	// 2. Достаем из БД
	camp, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err // Ошибку БД вернет сервис
	}

	// 3. Кладем в кеш асинхронно
	go func() {
		bytes, err := json.Marshal(camp)
		if err != nil {
			r.logger.Error().Err(err).Str("id", id).Msg("failed to marshal campaign for cache")
			return
		}

		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.cache.Set(bgCtx, key, bytes, r.ttl); err != nil {
			r.logger.Error().Err(err).Str(logkeys.CacheKey, key).Msg("failed to set campaign to cache")
		}
	}()

	return camp, nil
}

func (r *CampaignRepo) FindAll(ctx context.Context, filter domain.CampaignFilter) ([]*domain.Campaign, error) {
	return r.repo.FindAll(ctx, filter)
}

func (r *CampaignRepo) Count(ctx context.Context, filter domain.CampaignFilter) (int64, error) {
	return r.repo.Count(ctx, filter)
}

func (r *CampaignRepo) Delete(ctx context.Context, userID domain.UserID, id string) error {
	if err := r.repo.Delete(ctx, userID, id); err != nil {
		return err
	}
	r.invalidate(id)
	return nil
}

func (r *CampaignRepo) DeleteByUserID(ctx context.Context, userID domain.UserID) error {
	return r.repo.DeleteByUserID(ctx, userID)
}

func (r *CampaignRepo) invalidate(id string) {
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		key := r.buildKey(id)
		if err := r.cache.Delete(bgCtx, key); err != nil {
			r.logger.Error().Err(err).Str(logkeys.CacheKey, key).Msg("failed to delete campaign from cache")
		}
	}()
}
