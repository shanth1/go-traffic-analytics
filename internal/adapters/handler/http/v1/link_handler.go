package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shanth1/gotrace/internal/core/services"
	"github.com/shanth1/gotrace/internal/pkg/request"
)

type LinkHandler struct {
	service *services.LinkService
}

func NewLinkHandler(s *services.LinkService) *LinkHandler {
	return &LinkHandler{service: s}
}

// --- Campaigns ---

func (h *LinkHandler) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	campaigns, err := h.service.GetUserCampaigns(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": campaigns})
}

type createCampaignReq struct {
	Name string `json:"name"`
}

func (h *LinkHandler) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	var req createCampaignReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "bad request")
		return
	}

	camp, err := h.service.CreateCampaign(r.Context(), userID, req.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{"data": camp})
}

// --- Links ---

func (h *LinkHandler) GetLinksByCampaign(w http.ResponseWriter, r *http.Request) {
	campaignID := chi.URLParam(r, "id")

	links, err := h.service.GetLinks(r.Context(), campaignID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": links})
}

type createLinkReq struct {
	CampaignID string `json:"campaign_id"`
	TargetURL  string `json:"target_url"`
	CustomSlug string `json:"custom_slug"` // Optional
}

func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	userID := request.GetUserID(r)

	var req createLinkReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "bad request")
		return
	}

	link, err := h.service.CreateLink(r.Context(), userID, req.CampaignID, req.TargetURL, req.CustomSlug)
	if err != nil {
		if err == services.ErrLimitReached {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{"data": link})
}

func (h *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteLink(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
