package cached

import (
	"context"
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

func (r *CachedClickRepo) Save(_ context.Context, click *domain.ClickEvent) error {

}

func (r *CachedClickRepo) CountTotal(_ context.Context, filter ports.AnalyticsFilter) (int64, error) {

}

func (r *CachedClickRepo) GetTimeSeriesGrouped(_ context.Context, filter ports.AnalyticsFilter, dimension string, interval time.Duration) ([]domain.StackedPoint, error) {

}

func (r *CachedClickRepo) GetFlowData(_ context.Context, filter ports.AnalyticsFilter, stages []string) (*domain.SankeyData, error) {

}

func (r *CachedClickRepo) GetHeatmapData(_ context.Context, filter ports.AnalyticsFilter) ([]domain.HeatmapPoint, error) {

}

func (r *CachedClickRepo) GetTopStats(_ context.Context, filter ports.AnalyticsFilter, dimension string, limit int) ([]domain.CategoryStat, error) {

}
