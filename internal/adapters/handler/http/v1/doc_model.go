package v1

import "github.com/shanth1/gotrace/internal/core/domain"

// List Data Wrappers
type UserListData struct {
	Data []*domain.User `json:"data"`
}

type PlanListData struct {
	Data []*domain.Plan `json:"data"`
}

type CampaignListData struct {
	Data []*domain.Campaign `json:"data"`
}

type LinkListData struct {
	Data []*domain.Link `json:"data"`
}

// Single Item Data Wrappers
type CampaignData struct {
	Data *domain.Campaign `json:"data"`
}

type LinkData struct {
	Data *domain.Link `json:"data"`
}

// Analytics Data Wrappers
type StreamGraphData struct {
	Data []domain.StackedPoint `json:"data"`
}

type SankeyData struct {
	Data *domain.SankeyData `json:"data"`
}

type GeoData struct {
	Data interface{} `json:"data"`
}

type QualityData struct {
	Data *domain.TrafficQuality `json:"data"`
}

type HeatmapData struct {
	Data []domain.HeatmapPoint `json:"data"`
}

type StatsData struct {
	Data []domain.CategoryStat `json:"data"`
}

// Status Data Wrapper
type StatusData struct {
	Status string `json:"status"`
}

// Response Models
type UsersListResponse struct {
	Data UserListData `json:"data"`
}

type PlansListResponse struct {
	Data PlanListData `json:"data"`
}

type CampaignResponse struct {
	Data CampaignData `json:"data"`
}

type CampaignsListResponse struct {
	Data CampaignListData `json:"data"`
}

type LinkResponse struct {
	Data LinkData `json:"data"`
}

type LinksListResponse struct {
	Data LinkListData `json:"data"`
}

type LoginData struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

type LoginResponse struct {
	Data LoginData `json:"data"`
}

type CreatedUserData struct {
	Data *domain.User `json:"data"`
}

type CreatedUserResponse struct {
	Data *domain.User `json:"data"`
}

// Analytics Responses
type AnalyticsSummaryResponse struct {
	Data *domain.Summary `json:"data"`
}

type StreamGraphResponse struct {
	Data StreamGraphData `json:"data"`
}

type SankeyResponse struct {
	Data SankeyData `json:"data"`
}

type GeoResponse struct {
	Data GeoData `json:"data"`
}

type QualityResponse struct {
	Data QualityData `json:"data"`
}

type HeatmapResponse struct {
	Data HeatmapData `json:"data"`
}

type StatsResponse struct {
	Data StatsData `json:"data"`
}

type ProfileTreeResponse struct {
	Data *domain.HierarchyNode `json:"data"`
}

type StatusResponse struct {
	Data StatusData `json:"data"`
}
