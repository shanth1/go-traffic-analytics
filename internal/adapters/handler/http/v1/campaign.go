package v1

import (
	"net/http"
	"strconv"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/request"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
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
// @Summary Get campaigns
// @Description Get list of user campaigns with pagination
// @Tags Campaigns
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit (default 10)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {object} CampaignsListResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/campaigns [get]
func (h *CampaignHandler) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)
	q := r.URL.Query()

	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(q.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	filter := domain.CampaignFilter{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}

	campaigns, total, err := h.campSvc.GetCampaigns(r.Context(), filter)
	if err != nil {
		response.RespondWithError(w, r, err)
		return
	}

	if campaigns == nil {
		campaigns = []*domain.Campaign{}
	}

	response.Success(w, r, CampaignsListResponse{
		Data: campaigns,
		Meta: domain.PaginationMeta{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	})
}

// CreateCampaign godoc
// @Summary Create campaign
// @Description Create a new campaign
// @Tags Campaigns
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateCampaignRequest true "Campaign Name"
// @Success 201 {object} CampaignResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/campaigns [post]
func (h *CampaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	var req CreateCampaignRequest
	if err := request.DecodeJSON(w, r, &req); err != nil {
		response.ClientError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" {
		response.ClientError(w, r, http.StatusBadRequest, "name is required")
		return
	}

	camp, err := h.campSvc.CreateCampaign(r.Context(), userID, req.Name)
	if err != nil {
		response.RespondWithError(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().Str("campaign_id", string(camp.ID)).Str("user_id", string(userID)).Str("name", camp.Name).Msg("campaign_created")

	response.CreatedData(w, r, camp)
}
