package ports

import (
	"context"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
)

//go:generate mockgen -source=repositories.go -destination=mocks/repositories_mock.go -package=mocks

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error)
	IncrementClickCount(ctx context.Context, userID string) error
}

type CampaignRepository interface {
	Save(ctx context.Context, campaign *domain.Campaign) error
	FindByID(ctx context.Context, id string) (*domain.Campaign, error)
	FindAllByUserID(ctx context.Context, userID string) ([]*domain.Campaign, error)
	Delete(ctx context.Context, id string) error
}

type PlanRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Plan, error)
	FindDefault(ctx context.Context) (*domain.Plan, error) // Обычно "free"
	FindAll(ctx context.Context) ([]*domain.Plan, error)
	Save(ctx context.Context, plan *domain.Plan) error // Для админа/сидера
}

type LinkRepository interface {
	Save(ctx context.Context, link *domain.Link) error
	FindBySlug(ctx context.Context, slug string) (*domain.Link, error)
	FindAll(ctx context.Context) ([]*domain.Link, error)
	FindAllByCampaignID(ctx context.Context, campaignID string) ([]*domain.Link, error)
	CountByUserID(ctx context.Context, userID string) (int64, error)
}

type AnalyticsFilter struct {
	LinkID     string
	CampaignID string
	From       time.Time
	To         time.Time
}

type ClickRepository interface {
	Save(ctx context.Context, click *domain.ClickEvent) error

	CountTotal(ctx context.Context, filter AnalyticsFilter) (int64, error)

	// Для Streamgraph / Stacked Area Chart
	// dimension: "os", "browser", "country"
	// interval: 1h, 24h
	GetTimeSeriesGrouped(ctx context.Context, filter AnalyticsFilter, dimension string, interval time.Duration) ([]domain.StackedPoint, error)

	// Для Sankey Diagram
	// stages: []string{"referer", "device", "country"}
	GetFlowData(ctx context.Context, filter AnalyticsFilter, stages []string) (*domain.SankeyData, error)

	GetHeatmapData(ctx context.Context, filter AnalyticsFilter) ([]domain.HeatmapPoint, error)

	// Для Pie/Bar/Donut/Geo
	// dimension: "browser", "os", "country", "referer"
	// limit: например, топ 10 + "Others"
	GetTopStats(ctx context.Context, filter AnalyticsFilter, dimension string, limit int) ([]domain.CategoryStat, error)
}
