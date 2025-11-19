package ports

import (
	"context"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
)

type LinkRepository interface {
	Save(ctx context.Context, link *domain.Link) error
	FindBySlug(ctx context.Context, slug string) (*domain.Link, error)
	FindAll(ctx context.Context) ([]*domain.Link, error)
}

type AnalyticsFilter struct {
	LinkID     string
	CampaignID string
	From       time.Time
	To         time.Time
}

type ClickRepository interface {
	Save(ctx context.Context, click *domain.ClickEvent) error

	// Базовые метрики
	CountTotal(ctx context.Context, filter AnalyticsFilter) (int64, error)

	// Для Streamgraph / Stacked Area Chart
	// dimension: "os", "browser", "country"
	// interval: 1h, 24h
	GetTimeSeriesGrouped(ctx context.Context, filter AnalyticsFilter, dimension string, interval time.Duration) ([]domain.StackedPoint, error)

	// Для Sankey Diagram
	// stages: []string{"referer", "device", "country"}
	GetFlowData(ctx context.Context, filter AnalyticsFilter, stages []string) (*domain.SankeyData, error)
}
