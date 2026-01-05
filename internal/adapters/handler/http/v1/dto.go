package v1

import (
	"errors"
	"strings"

	"github.com/shanth1/gotrace/internal/core/domain"
)

// --- Requests ---

type CreateLinkRequest struct {
	CampaignID string `json:"campaign_id"`
	TargetURL  string `json:"target_url" binding:"required" minLength:"1" example:"https://google.com"`
}

func (r CreateLinkRequest) Validate() error {
	if strings.TrimSpace(r.TargetURL) == "" {
		return errors.New("target_url is required")
	}
	return nil
}

type CreateCampaignRequest struct {
	Name string `json:"name" binding:"required" minLength:"1" example:"My Campaign"`
}

func (r CreateCampaignRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}
	return nil
}

type RegisterReq struct {
	Email    string `json:"email" binding:"required" minLength:"5" example:"user@example.com"`
	Password string `json:"password" binding:"required" minLength:"6" example:"secret123"`
}

func (r RegisterReq) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email is required")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

type LoginReq struct {
	Email    string `json:"email" binding:"required" minLength:"1" example:"user@example.com"`
	Password string `json:"password" binding:"required" minLength:"1" example:"secret123"`
}

func (r LoginReq) Validate() error {
	if strings.TrimSpace(r.Email) == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type UpdateUserStatusReq struct {
	IsActive bool `json:"is_active" example:"true"`
}

func (r UpdateUserStatusReq) Validate() error {
	// Булево значение всегда присутствует (true или false),
	// дополнительная валидация не требуется, но метод нужен для унификации.
	return nil
}

type UpdateUserPlanReq struct {
	PlanID string `json:"plan_id" binding:"required" minLength:"1" example:"plan_pro"`
}

func (r UpdateUserPlanReq) Validate() error {
	if strings.TrimSpace(r.PlanID) == "" {
		return errors.New("plan_id is required")
	}
	return nil
}

// --- Responses (Data Wrappers) ---

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

type AnalyticsSummaryResponse struct {
	Data domain.Summary `json:"data"`
}

type StreamGraphResponse struct {
	Data []map[string]any `json:"data"`
}

type SankeyResponse struct {
	Data domain.SankeyData `json:"data"`
}

type GeoResponse struct { // Исправлена опечатка GeoRespons -> GeoResponse
	Data []domain.GeoStat `json:"data"`
}

type QualityResponse struct {
	Data domain.TrafficQuality `json:"data"`
}

type HeatmapResponse struct {
	Data []domain.HeatmapPoint `json:"data"`
}

type StatsResponse struct {
	Data []domain.CategoryStat `json:"data"`
}
