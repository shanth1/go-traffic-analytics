package ports

import (
	"context"

	"github.com/shanth1/gotrace/internal/core/domain"
)

//go:generate mockgen -source=services.go -destination=mocks/service_mock.go -package=mocks

type AuthService interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, *domain.User, error)
}

type AnalyticsService interface {
	GetSummary(ctx context.Context, filter AnalyticsFilter) (*domain.Summary, error)
	GetStreamGraphData(ctx context.Context, filter AnalyticsFilter, groupBy string) ([]domain.StackedPoint, error)
	GetSankeyData(ctx context.Context, filter AnalyticsFilter, stages []string) (*domain.SankeyData, error)
	GetGeoDistribution(ctx context.Context, filter AnalyticsFilter) ([]domain.GeoStat, error)
	GetTrafficQuality(ctx context.Context, filter AnalyticsFilter) (*domain.TrafficQuality, error)
	GetHeatmapData(ctx context.Context, filter AnalyticsFilter) ([]domain.HeatmapPoint, error)
	GetCategoryStats(ctx context.Context, filter AnalyticsFilter, dimension string) ([]domain.CategoryStat, error)
}

type LinkService interface {
	CreateLink(ctx context.Context, userID domain.UserID, campaignID, targetURL, customSlug string) (*domain.Link, error)
	GetLinks(ctx context.Context, campaignID string) ([]*domain.Link, error)
	GetUserCampaigns(ctx context.Context, userID domain.UserID) ([]*domain.Campaign, error)
	CreateCampaign(ctx context.Context, userID domain.UserID, name string) (*domain.Campaign, error)
	DeleteLink(ctx context.Context, id domain.LinkID) error
	GetUserHierarchy(ctx context.Context, userID domain.UserID) (*domain.HierarchyNode, error)
}

type RedirectService interface {
	ProcessRedirect(ctx context.Context, slug, ip, userAgentString, referer string) (string, error)
}

type UserService interface {
	GetAllUsers(ctx context.Context, page, limit int) ([]*domain.User, error)
	SetUserStatus(ctx context.Context, userID domain.UserID, isActive bool) error
	ChangeUserPlan(ctx context.Context, userID domain.UserID, planID string) error
}

type BillingService interface {
	GetAllPlans(ctx context.Context) ([]*domain.Plan, error)
}
