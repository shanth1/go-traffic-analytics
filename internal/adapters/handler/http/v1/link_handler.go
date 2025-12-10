package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (h *LinkHandler) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	campaigns, err := h.service.GetUserCampaigns(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{"data": campaigns})
}

type createCampaignReq struct {
	Name string `json:"name"`
}

func (h *LinkHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	var req createCampaignReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	camp, err := h.service.CreateCampaign(r.Context(), userID, req.Name)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{"data": camp})
}

// --- Links ---

func (h *LinkHandler) GetLinksByCampaign(w http.ResponseWriter, r *http.Request) {
	campaignID := chi.URLParam(r, "id")

	links, err := h.service.GetLinks(r.Context(), campaignID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{"data": links})
}

func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	// TODO: validation
	if req.TargetURL == "" {
		response.JSON(w, http.StatusBadRequest, "target_url is required")
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

	response.JSON(w, http.StatusCreated, map[string]interface{}{"data": link})
}

func (h *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteLink(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
