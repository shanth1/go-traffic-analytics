package ports

import (
	"context"
	"io"

	"github.com/shanth1/gotrace/internal/core/domain"
)

//go:generate mockgen -source=services.go -destination=mocks/service_mock.go -package=mocks

type AuthService interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, *domain.User, error)
}

type AnalyticsService interface {
	GetSummary(ctx context.Context, filter domain.AnalyticsFilter) (*domain.Summary, error)
	GetStreamGraphData(ctx context.Context, filter domain.AnalyticsFilter, groupBy string) ([]domain.StackedPoint, error)
	GetSankeyData(ctx context.Context, filter domain.AnalyticsFilter, stages []string) (*domain.SankeyData, error)
	GetGeoStats(ctx context.Context, filter domain.AnalyticsFilter) (domain.GeoStats, error)
	GetTrafficQuality(ctx context.Context, filter domain.AnalyticsFilter) (*domain.TrafficQuality, error)
	GetHeatmapData(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.HeatmapPoint, error)
	GetCategoryStats(ctx context.Context, filter domain.AnalyticsFilter, dimension string) ([]domain.CategoryStat, error)
	ExportData(ctx context.Context, filter domain.AnalyticsFilter) (io.Reader, string, error)
}

type LinkService interface {
	CreateLink(ctx context.Context, cmd domain.CreateLinkCmd) (*domain.Link, error)

	// TODO: ?
	GetLinkList(ctx context.Context, filter domain.LinkFilter) ([]*domain.Link, int64, error) // Returns total count

	DeleteLink(ctx context.Context, userID domain.UserID, id domain.LinkID) error
}

type CampaignService interface {
	CreateCampaign(ctx context.Context, userID domain.UserID, name string) (*domain.Campaign, error)

	// TODO: ?
	GetCampaigns(ctx context.Context, filter domain.CampaignFilter) ([]*domain.Campaign, int64, error) // Returns total count

	DeleteCampaign(ctx context.Context, userID domain.UserID, id string) error
}

type RedirectService interface {
	Process(ctx context.Context, slug string, meta domain.RequestMetadata) (string, error)
}

type UserService interface {
	GetAll(ctx context.Context, page, limit int) ([]*domain.User, int64, error) // Returns total count
	SetStatus(ctx context.Context, userID domain.UserID, isActive bool) error
	ChangePlan(ctx context.Context, userID domain.UserID, planID string) error
	GetHierarchy(ctx context.Context, id domain.UserID) (*domain.HierarchyNode, error)
	DeleteUser(ctx context.Context, userID domain.UserID) error
}

type BillingService interface {
	GetAllPlans(ctx context.Context) ([]*domain.Plan, error)
}

type FeedbackService interface {
	SendFeedback(ctx context.Context, cmd domain.SendFeedbackCmd) error
}
