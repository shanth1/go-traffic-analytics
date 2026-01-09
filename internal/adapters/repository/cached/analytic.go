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

type AnalyticRepo struct {
	repo  ports.AnalyticsRepository
	cache ports.Cache
	ttl   time.Duration
}

func NewAnalyicRepo(repo ports.AnalyticsRepository, cache ports.Cache, ttl time.Duration) ports.AnalyticsRepository {
	return &AnalyticRepo{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (r *AnalyticRepo) buildKey(prefix string, params ...interface{}) string {
	key := fmt.Sprintf("gotrace:analytics:%s", prefix)
	for _, p := range params {
		b, _ := json.Marshal(p)
		key += ":" + string(b)
	}
	return key
}

func (r *AnalyticRepo) SaveBatch(ctx context.Context, events []*domain.ClickEvent) error {
	const op = "cached.AnalyticRepo.SaveBatch"

	if err := r.repo.SaveBatch(ctx, events); err != nil {
		return ops.E(op, err)
	}

	return nil
}

func (r *AnalyticRepo) CountTotal(ctx context.Context, filter domain.AnalyticsFilter) (int64, error) {
	const op = "cached.AnalyticRepo.CountTotal"

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
	}
	// else if !errors.Is(err, ports.ErrCacheMiss) {
	// TODO: logging (cache error)
	// }

	count, err := r.repo.CountTotal(ctx, filter)
	if err != nil {
		return 0, ops.E(op, err)
	}

	go func() {
		bytes, _ := json.Marshal(count)
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return count, nil
}

func (r *AnalyticRepo) GetTimeSeriesGrouped(ctx context.Context, filter domain.AnalyticsFilter, dimension string, interval time.Duration) ([]domain.StackedPoint, error) {
	const op = "cached.AnalyticRepo.GetTimeSeriesGrouped"

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
	}
	// else if !errors.Is(err, ports.ErrCacheMiss) {
	// TODO: logging (cache error)
	// }

	result, err := r.repo.GetTimeSeriesGrouped(ctx, filter, dimension, interval)
	if err != nil {
		return nil, ops.E(op, err)
	}

	go func() {
		bytes, _ := json.Marshal(result)
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return result, nil
}

func (r *AnalyticRepo) GetFlowData(ctx context.Context, filter domain.AnalyticsFilter, stages []string) (*domain.SankeyData, error) {
	const op = "cached.AnalyticRepo.GetFlowData"

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
	}
	// else if !errors.Is(err, ports.ErrCacheMiss) {
	// TODO: logging (cache error)
	// }

	data, err := r.repo.GetFlowData(ctx, filter, stages)
	if err != nil {
		return nil, ops.E(op, err)
	}

	go func() {
		bytes, _ := json.Marshal(data)
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return data, nil
}

func (r *AnalyticRepo) GetHeatmapData(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.HeatmapPoint, error) {
	const op = "cached.AnalyticRepo.GetHeatmapData"

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
	}
	// else if !errors.Is(err, ports.ErrCacheMiss) {
	// TODO: logging (cache error)
	// }

	points, err := r.repo.GetHeatmapData(ctx, filter)
	if err != nil {
		return nil, ops.E(op, err)
	}

	go func() {
		bytes, _ := json.Marshal(points)
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return points, nil
}

func (r *AnalyticRepo) GetTopStats(ctx context.Context, filter domain.AnalyticsFilter, dimension string, limit int) ([]domain.CategoryStat, error) {
	const op = "cached.AnalyticRepo.GetTopStats"

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
	}
	// else if !errors.Is(err, ports.ErrCacheMiss) {
	// TODO: logging (cache error)
	// }

	stats, err := r.repo.GetTopStats(ctx, filter, dimension, limit)
	if err != nil {
		return nil, ops.E(op, err)
	}

	go func() {
		bytes, _ := json.Marshal(stats)
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return stats, nil
}
