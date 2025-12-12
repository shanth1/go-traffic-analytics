package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

type AdminHandler struct {
	userService ports.UserService
	logger      log.Logger
}

func NewAdminHandler(u ports.UserService, l log.Logger) *AdminHandler {
	return &AdminHandler{
		userService: u,
		logger:      l,
	}
}

// --- Handlers ---

// GetUsers godoc
// @Summary      List users
// @Description  Get paginated users list (Admin)
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Param        page query int false "Page number"
// @Success      200  {object}  UsersListResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/api/v1/admin/users [get]
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

	response.JSON(w, http.StatusOK, UsersListResponse{Data: users})
}

type UpdateUserStatusReq struct {
	IsActive bool `json:"is_active"`
}

// UpdateUserStatus godoc
// @Summary      Update user status
// @Description  Activate/Deactivate user
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id      path string              true "User ID"
// @Param        request body UpdateUserStatusReq true "Status"
// @Success      200  {object}  map[string]string      "Status: updated"
// @Failure      400  {object}  response.ErrorResponse
// @Router       /api/v1/admin/users/{id}/status [patch]
func (h *AdminHandler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateUserStatusReq
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

type UpdateUserPlanReq struct {
	PlanID string `json:"plan_id"`
}

// UpdateUserPlan godoc
// @Summary      Change user plan
// @Description  Set new plan for user
// @Tags         Admin
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id      path string            true "User ID"
// @Param        request body UpdateUserPlanReq true "New Plan ID"
// @Success      200  {object}  map[string]string      "Status: updated"
// @Failure      400  {object}  response.ErrorResponse
// @Router       /api/v1/admin/users/{id}/plan [patch]
func (h *AdminHandler) UpdateUserPlan(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateUserPlanReq
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

// GetPlans godoc
// @Summary      List plans
// @Description  Get all available plans
// @Tags         Admin
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  PlansListResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/admin/plans [get]
func (h *AdminHandler) GetPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.userService.GetAllPlans(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{"data": plans})
}
