package v1

import (
	"encoding/json"
	"net/http"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
	"github.com/shanth1/gotrace/internal/pkg/request"
)

type CampaignHandler struct {
	campSvc ports.CampaignService
}

func NewCampaignHandler(cs ports.CampaignService) *CampaignHandler {
	return &CampaignHandler{
		campSvc: cs,
	}
}

// GetCampaigns godoc
// @Summary      Get campaigns
// @Description  Get list of user campaigns
// @Tags         Campaigns
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  CampaignsListResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/campaigns [get]
func (h *CampaignHandler) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	campaigns, err := h.campSvc.GetCampaigns(r.Context(), userID)
	if err != nil {
		response.ServerError(w, r, err)
		return
	}

	if campaigns == nil {
		campaigns = []*domain.Campaign{}
	}

	response.JSON(w, http.StatusOK, response.Envelope{"data": campaigns})
}

// CreateCampaign godoc
// @Summary      Create campaign
// @Description  Create a new campaign
// @Tags         Campaigns
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body CreateCampaignRequest true "Campaign Name"
// @Success      201  {object}  CampaignResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/campaigns [post]
func (h *CampaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	var req CreateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	camp, err := h.campSvc.CreateCampaign(r.Context(), userID, req.Name)
	if err != nil {
		response.ServerError(w, r, err)
		return
	}

	response.JSON(w, http.StatusCreated, response.Envelope{"data": camp})
}
