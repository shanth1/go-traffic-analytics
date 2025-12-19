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

type CachedClickRepo struct {
	repo  ports.ClickRepository
	cache ports.Cache
	ttl   time.Duration
}

func NewClickRepo(repo ports.ClickRepository, cache ports.Cache, ttl time.Duration) ports.ClickRepository {
	return &CachedClickRepo{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (r *CachedClickRepo) buildKey(prefix string, params ...interface{}) string {
	key := fmt.Sprintf("gotrace:analytics:%s", prefix)
	for _, p := range params {
		b, _ := json.Marshal(p)
		key += ":" + string(b)
	}
	return key
}

func (r *CachedClickRepo) Save(ctx context.Context, click *domain.ClickEvent) error {
	return r.repo.Save(ctx, click)
}

func (r *CachedClickRepo) CountTotal(ctx context.Context, filter ports.AnalyticsFilter) (int64, error) {
	key := r.buildKey("count", filter)

	val, err := r.cache.Get(ctx, key)
	if err == nil {
		if bytesVal, ok := val.([]byte); ok {
			var count int64
			if jsonErr := json.Unmarshal(bytesVal, &count); jsonErr == nil {
				return count, nil
			}
			// TODO: logging (unmarshal error)
		}
	} else if !errors.Is(err, ports.ErrCacheMiss) {
		// TODO: logging (cache error)
	}

	count, err := r.repo.CountTotal(ctx, filter)
	if err != nil {
		return 0, err
	}

	go func() {
		bytes, _ := json.Marshal(count)
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return count, nil
}

func (r *CachedClickRepo) GetTimeSeriesGrouped(ctx context.Context, filter ports.AnalyticsFilter, dimension string, interval time.Duration) ([]domain.StackedPoint, error) {
	key := r.buildKey("timeseries", filter, dimension, interval)

	val, err := r.cache.Get(ctx, key)
	if err == nil {
		if bytesVal, ok := val.([]byte); ok {
			var result []domain.StackedPoint
			if jsonErr := json.Unmarshal(bytesVal, &result); jsonErr == nil {
				return result, nil
			}
			// TODO: logging (unmarshal error)
		}
	} else if !errors.Is(err, ports.ErrCacheMiss) {
		// TODO: logging (cache error)
	}

	result, err := r.repo.GetTimeSeriesGrouped(ctx, filter, dimension, interval)
	if err != nil {
		return nil, err
	}

	go func() {
		bytes, _ := json.Marshal(result)
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return result, nil
}

func (r *CachedClickRepo) GetFlowData(ctx context.Context, filter ports.AnalyticsFilter, stages []string) (*domain.SankeyData, error) {
	key := r.buildKey("flow", filter, stages)

	val, err := r.cache.Get(ctx, key)
	if err == nil {
		if bytesVal, ok := val.([]byte); ok {
			var data domain.SankeyData
			if jsonErr := json.Unmarshal(bytesVal, &data); jsonErr == nil {
				return &data, nil
			}
			// TODO: logging (unmarshal error)
		}
	} else if !errors.Is(err, ports.ErrCacheMiss) {
		// TODO: logging (cache error)
	}

	data, err := r.repo.GetFlowData(ctx, filter, stages)
	if err != nil {
		return nil, err
	}

	go func() {
		bytes, _ := json.Marshal(data)
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return data, nil
}

func (r *CachedClickRepo) GetHeatmapData(ctx context.Context, filter ports.AnalyticsFilter) ([]domain.HeatmapPoint, error) {
	key := r.buildKey("heatmap", filter)

	val, err := r.cache.Get(ctx, key)
	if err == nil {
		if bytesVal, ok := val.([]byte); ok {
			var points []domain.HeatmapPoint
			if jsonErr := json.Unmarshal(bytesVal, &points); jsonErr == nil {
				return points, nil
			}
			// TODO: logging (unmarshal error)
		}
	} else if !errors.Is(err, ports.ErrCacheMiss) {
		// TODO: logging (cache error)
	}

	points, err := r.repo.GetHeatmapData(ctx, filter)
	if err != nil {
		return nil, err
	}

	go func() {
		bytes, _ := json.Marshal(points)
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return points, nil
}

func (r *CachedClickRepo) GetTopStats(ctx context.Context, filter ports.AnalyticsFilter, dimension string, limit int) ([]domain.CategoryStat, error) {
	key := r.buildKey("topstats", filter, dimension, limit)

	val, err := r.cache.Get(ctx, key)
	if err == nil {
		if bytesVal, ok := val.([]byte); ok {
			var stats []domain.CategoryStat
			if jsonErr := json.Unmarshal(bytesVal, &stats); jsonErr == nil {
				return stats, nil
			}
			// TODO: logging (unmarshal error)
		}
	} else if !errors.Is(err, ports.ErrCacheMiss) {
		// TODO: logging (cache error)
	}

	stats, err := r.repo.GetTopStats(ctx, filter, dimension, limit)
	if err != nil {
		return nil, err
	}

	go func() {
		bytes, _ := json.Marshal(stats)
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return stats, nil
}
