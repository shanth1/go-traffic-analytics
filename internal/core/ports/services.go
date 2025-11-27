package ports

import (
	"context"

	"github.com/shanth1/gotrace/internal/core/domain"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, *domain.User, error)
}

type LinkService interface {
	CreateLink(ctx context.Context, userID, campaignID, targetURL, customSlug string) (*domain.Link, error)
	GetLinks(ctx context.Context, campaignID string) ([]*domain.Link, error)
	GetUserCampaigns(ctx context.Context, userID string) ([]*domain.Campaign, error)
	CreateCampaign(ctx context.Context, userID, name string) (*domain.Campaign, error)
	DeleteLink(ctx context.Context, id string) error
}

type AnalyticsService interface {
	GetSummary(ctx context.Context, filter AnalyticsFilter) (map[string]interface{}, error)
	GetStreamGraphData(ctx context.Context, filter AnalyticsFilter, groupBy string) ([]domain.StackedPoint, error)
	GetSankeyData(ctx context.Context, filter AnalyticsFilter, stages []string) (*domain.SankeyData, error)
	GetGeoDistribution(ctx context.Context, filter AnalyticsFilter) ([]domain.TimeSeriesPoint, error) // Упрощено для примера
	GetTrafficQuality(ctx context.Context, filter AnalyticsFilter) (map[string]int, error)
}
