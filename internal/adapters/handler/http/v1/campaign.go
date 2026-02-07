package v1

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotools/ops"
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
// @Failure 401 {object} response.ErrorWrapper "Unauthorized"
// @Failure 500 {object} response.ErrorWrapper
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
		response.Error(w, r, err)
		return
	}

	if campaigns == nil {
		campaigns = []*domain.Campaign{}
	}

	log.FromContext(r.Context()).Info().
		Int("count", len(campaigns)).
		Int("offset", offset).
		Int("limit", limit).
		Int64("total", total).
		Msg("user_campaigns_listed")

	response.JSON(w, r, http.StatusOK, CampaignsListResponse{
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
// @Failure 400 {object} response.ErrorWrapper
// @Failure 401 {object} response.ErrorWrapper "Unauthorized"
// @Failure 500 {object} response.ErrorWrapper
// @Router /api/v1/campaigns [post]
func (h *CampaignHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	const op = "v1.CampaignHandler.CreateCampaign"

	userID := request.GetUserID(r)

	var req CreateCampaignRequest
	if err := request.DecodeJSON(w, r, &req); err != nil {
		err := ops.WrapMsg(op, ops.KindInvalid, err, "invalid request body")
		response.Error(w, r, err)
		return
	}

	camp, err := h.campSvc.CreateCampaign(r.Context(), userID, req.Name)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().
		Str("campaign_id", camp.ID).
		Str("name", camp.Name).
		Msg("campaign_created")

	response.Created(w, r, camp)
}

// DeleteCampaign godoc
// @Summary Delete campaign
// @Description Delete a campaign and cascade soft-delete all its links
// @Tags Campaigns
// @Security BearerAuth
// @Param id path string true "Campaign ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} response.ErrorWrapper "Invalid ID"
// @Failure 401 {object} response.ErrorWrapper "Unauthorized"
// @Failure 404 {object} response.ErrorWrapper "Not Found"
// @Failure 500 {object} response.ErrorWrapper
// @Router /api/v1/campaigns/{id} [delete]
func (h *CampaignHandler) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	const op = "v1.CampaignHandler.DeleteCampaign"

	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, r, ops.WrapMsg(op, ops.KindInvalid, nil, "campaign id is required"))
		return
	}

	userID := request.GetUserID(r)

	if err := h.campSvc.DeleteCampaign(r.Context(), userID, id); err != nil {
		response.Error(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().
		Str("campaign_id", id).
		Msg("campaign_deleted")

	response.NoContent(w, r)
}
