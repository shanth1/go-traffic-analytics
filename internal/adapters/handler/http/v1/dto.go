package v1

import "github.com/shanth1/gotrace/internal/core/domain"

// --- Requests ---

type CreateLinkRequest struct {
	CampaignID string `json:"campaign_id"`
	TargetURL  string `json:"target_url" binding:"required" example:"https://google.com"`
}

type CreateCampaignRequest struct {
	Name string `json:"name" binding:"required" example:"My Campaign"`
}

type RegisterReq struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

type UpdateUserStatusReq struct {
	IsActive bool `json:"is_active" example:"true"`
}

type UpdateUserPlanReq struct {
	PlanID string `json:"plan_id" binding:"required" example:"plan_pro"`
}

// --- Responses (Data Wrappers) ---

// StatusData represents a generic status message inside data
type StatusData struct {
	Status string `json:"status" example:"updated"`
}

type StatusResponse struct {
	Data StatusData `json:"data"`
}

// Auth
type LoginData struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

type LoginResponse struct {
	Data LoginData `json:"data"`
}

type RegisterResponse struct {
	Data *domain.User `json:"data"`
}

// Users
type UsersListResponse struct {
	Data []*domain.User `json:"data"`
}

type ProfileTreeResponse struct {
	// Предпологаем, что HierarchyNode определен в domain
	Data domain.HierarchyNode `json:"data"`
}

// Plans
type PlansListResponse struct {
	Data []*domain.Plan `json:"data"`
}

// Campaigns
type CampaignsListResponse struct {
	Data []*domain.Campaign `json:"data"`
}

type CampaignResponse struct {
	Data *domain.Campaign `json:"data"`
}

// Links
type LinksListResponse struct {
	Data []*domain.Link `json:"data"`
}

type LinkResponse struct {
	Data *domain.Link `json:"data"`
}

// Analytics Wrappers
// Используем any, если точные типы в domain сложны для импорта здесь,
// но лучше использовать точные типы из domain, если они экспортируемы.

type AnalyticsSummaryResponse struct {
	Data domain.Summary `json:"data"`
}

type StreamGraphResponse struct {
	// Предполагаем тип возвращаемого значения сервиса
	Data []map[string]any `json:"data"`
}

type SankeyResponse struct {
	Data domain.SankeyData `json:"data"`
}

type GeoResponse struct {
	Data []domain.GeoStat `json:"data"`
}

type QualityResponse struct {
	Data []domain.TrafficQuality `json:"data"`
}

type HeatmapResponse struct {
	Data []domain.HeatmapPoint `json:"data"`
}

type StatsResponse struct {
	Data []domain.CategoryStat `json:"data"`
}
