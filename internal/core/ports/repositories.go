package ports

import (
	"context"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
)

//go:generate mockgen -source=repositories.go -destination=mocks/repositories_mock.go -package=mocks

type Transactor interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
}

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id domain.UserID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error)
	Count(ctx context.Context) (int64, error)

	IncrementUsage(ctx context.Context, userID domain.UserID, delta int) error
	ResetUsage(ctx context.Context, userID domain.UserID) error

	Delete(ctx context.Context, id domain.UserID) error
}

type CampaignRepository interface {
	Save(ctx context.Context, campaign *domain.Campaign) error
	FindByID(ctx context.Context, id string) (*domain.Campaign, error)
	FindAll(ctx context.Context, filter domain.CampaignFilter) ([]*domain.Campaign, error)
	Count(ctx context.Context, filter domain.CampaignFilter) (int64, error)
	Delete(ctx context.Context, userID domain.UserID, id string) error
	DeleteByUserID(ctx context.Context, userID domain.UserID) error
}

type PlanRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Plan, error)
	FindDefault(ctx context.Context) (*domain.Plan, error)
	FindAll(ctx context.Context) ([]*domain.Plan, error)
	Save(ctx context.Context, plan *domain.Plan) error
}

type LinkRepository interface {
	Save(ctx context.Context, link *domain.Link) error
	FindBySlug(ctx context.Context, slug string) (*domain.Link, error)
	FindByID(ctx context.Context, id domain.LinkID) (*domain.Link, error)
	FindAll(ctx context.Context, filter domain.LinkFilter) ([]*domain.Link, error)
	Count(ctx context.Context, filter domain.LinkFilter) (int64, error)
	Delete(ctx context.Context, userID domain.UserID, id domain.LinkID) error
	DeleteByUserID(ctx context.Context, userID domain.UserID) error
	DeleteByCampaignID(ctx context.Context, campaignID string) error
}

type GeoProvider interface {
	Lookup(ctx context.Context, ip string) (*domain.GeoLocation, error)
}

// --- ANALYTICS / CLICK PROCESSING ---

type EventIngestor interface {
	TrackClick(ctx context.Context, event *domain.ClickEvent) error
	Close() error
}

type AnalyticsRepository interface {
	SaveBatch(ctx context.Context, events []*domain.ClickEvent) error
	CountTotal(ctx context.Context, filter domain.AnalyticsFilter) (int64, error)

	// Streamgraph / Stacked Area Chart
	// dimension: "os", "browser", "country"
	// interval: 1h, 24h
	GetTimeSeriesGrouped(ctx context.Context, filter domain.AnalyticsFilter, dimension string, interval time.Duration) ([]domain.StackedPoint, error)

	// Sankey Diagram
	// stages: []string{"referer", "device", "country"}
	GetFlowData(ctx context.Context, filter domain.AnalyticsFilter, stages []string) (*domain.SankeyData, error)

	GetHeatmapData(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.HeatmapPoint, error)

	// Pie/Bar/Donut/Geo
	// dimension: "browser", "os", "country", "referer"
	GetTopStats(ctx context.Context, filter domain.AnalyticsFilter, dimension string, limit int) ([]domain.CategoryStat, error)

	GetRawEvents(ctx context.Context, filter domain.AnalyticsFilter) ([]*domain.ClickEvent, error)

	// TODO:
	// AnonymizeUserData(ctx context.Context, userID domain.UserID) error
}
