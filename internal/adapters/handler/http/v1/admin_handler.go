package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

type AdminHandler struct {
	userService ports.UserService
}

func NewAdminHandler(u ports.UserService) *AdminHandler {
	return &AdminHandler{userService: u}
}

// --- Handlers ---

func (h *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit := 20

	users, err := h.userService.GetAllUsers(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{"data": users})
}

type updateUserStatusReq struct {
	IsActive bool `json:"is_active"`
}

func (h *AdminHandler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateUserStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	if err := h.userService.SetUserStatus(r.Context(), id, req.IsActive); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

type updateUserPlanReq struct {
	PlanID string `json:"plan_id"`
}

func (h *AdminHandler) UpdateUserPlan(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateUserPlanReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	if err := h.userService.ChangeUserPlan(r.Context(), id, req.PlanID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *AdminHandler) GetPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.userService.GetAllPlans(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{"data": plans})
}
