package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/core/services"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
	"github.com/shanth1/gotrace/internal/pkg/request"
)

type LinkHandler struct {
	service ports.LinkService
}

func NewLinkHandler(s ports.LinkService) *LinkHandler {
	return &LinkHandler{service: s}
}

// --- Campaigns ---

// GetCampaigns godoc
// @Summary      Get campaigns
// @Description  Get list of user campaigns
// @Tags         Campaigns
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  CampaignsListResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/campaigns [get]
func (h *LinkHandler) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	campaigns, err := h.service.GetUserCampaigns(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if campaigns == nil {
		campaigns = []*domain.Campaign{}
	}

	response.JSON(w, http.StatusOK, CampaignsListResponse{Data: campaigns})
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
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      400  {object}  response.ErrorResponse
// @Router       /api/v1/campaigns [post]
func (h *LinkHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	var req CreateCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	camp, err := h.service.CreateCampaign(r.Context(), userID, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, CampaignResponse{Data: camp})
}

// --- Links ---

// GetLinksByCampaign godoc
// @Summary      Get links
// @Description  Get all links for a specific campaign
// @Tags         Links
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Campaign ID"
// @Success      200  {object}  LinksListResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/campaigns/{id}/links [get]
func (h *LinkHandler) GetLinksByCampaign(w http.ResponseWriter, r *http.Request) {
	campaignID := chi.URLParam(r, "id")

	links, err := h.service.GetLinks(r.Context(), campaignID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if links == nil {
		links = []*domain.Link{}
	}

	response.JSON(w, http.StatusOK, LinksListResponse{Data: links})
}

// CreateLink godoc
// @Summary      Create link
// @Description  Shorten a URL
// @Tags         Links
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body CreateLinkRequest true "Link Info"
// @Success      201  {object}  LinkResponse
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      403  {object}  response.ErrorResponse "Limit reached"
// @Router       /api/v1/links [post]
func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	// TODO: validation
	if req.TargetURL == "" {
		response.Error(w, http.StatusBadRequest, "target_url is required")
		return
	}

	link, err := h.service.CreateLink(r.Context(), userID, req.CampaignID, req.TargetURL, req.CustomSlug)
	if err != nil {
		if err == services.ErrLimitReached {
			response.Error(w, http.StatusForbidden, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, LinkResponse{Data: link})
}

// DeleteLink godoc
// @Summary      Delete link
// @Description  Remove a link
// @Tags         Links
// @Security     BearerAuth
// @Param        id   path      string  true  "Link ID"
// @Success      204  {string}  string  "No Content"
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/links/{id} [delete]
func (h *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteLink(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Del("Content-Type")
	w.WriteHeader(http.StatusNoContent)
}

// GetProfileTree godoc
// @Summary      Profile Hierarchy
// @Description  Get hierarchical structure of User -> Campaigns -> Links for Tree visualization
// @Tags         Campaigns
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  domain.HierarchyNode
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/campaigns/tree [get]
func (h *LinkHandler) GetProfileTree(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	tree, err := h.service.GetUserHierarchy(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{"data": tree})
}
