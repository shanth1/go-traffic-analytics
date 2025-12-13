package v1

import "github.com/shanth1/gotrace/internal/core/domain"

type UsersListResponse struct {
	Data []*domain.User `json:"data"`
}

type PlansListResponse struct {
	Data []*domain.Plan `json:"data"`
}

type CampaignResponse struct {
	Data *domain.Campaign `json:"data"`
}

type CampaignsListResponse struct {
	Data []*domain.Campaign `json:"data"`
}

type LinkResponse struct {
	Data *domain.Link `json:"data"`
}

type LinksListResponse struct {
	Data []*domain.Link `json:"data"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

// Analytics Wrappers

type AnalyticsSummaryResponse struct {
	Data map[string]interface{} `json:"data"`
}

type StreamGraphResponse struct {
	Data []domain.StackedPoint `json:"data"`
}

type SankeyResponse struct {
	Data *domain.SankeyData `json:"data"`
}

type GeoResponse struct {
	Data interface{} `json:"data"`
}

type QualityResponse struct {
	Data map[string]int `json:"data"`
}

type HeatmapResponse struct {
	Data []domain.HeatmapPoint `json:"data"`
}

type StatsResponse struct {
	Data []domain.CategoryStat `json:"data"`
}

type ProfileTreeResponse struct {
	Data *domain.HierarchyNode `json:"data"`
}
