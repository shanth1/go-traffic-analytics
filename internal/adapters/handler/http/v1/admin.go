package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

type AdminHandler struct {
	userService ports.UserService
}

func NewAdminHandler(u ports.UserService) *AdminHandler {
	return &AdminHandler{
		userService: u,
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
// @Failure      401  {object}  response.ErrorResponse "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse "Forbidden"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/admin/users [get]
func (h *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit := 20

	users, err := h.userService.GetAll(r.Context(), page, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if users == nil {
		users = []*domain.User{}
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
// @Param        id       path domain.UserID                true "User ID"
// @Param        request  body UpdateUserStatusReq true     "Status"
// @Success      200      {object}  map[string]string       "Status: updated"
// @Failure      400      {object}  response.ErrorResponse
// @Failure      401      {object}  response.ErrorResponse  "Unauthorized"
// @Failure      403      {object}  response.ErrorResponse  "Forbidden"
// @Router       /api/v1/admin/users/{id}/status [patch]
func (h *AdminHandler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(chi.URLParam(r, "id"))

	var req UpdateUserStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	if err := h.userService.SetStatus(r.Context(), id, req.IsActive); err != nil {
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
// @Param        id      path   domain.UserID            true "User ID"
// @Param        request body   UpdateUserPlanReq true  "New Plan ID"
// @Success      200  {object}  map[string]string       "Status: updated"
// @Failure      401  {object}  response.ErrorResponse  "Unauthorized"
// @Failure      403  {object}  response.ErrorResponse  "Forbidden"
// @Failure      400  {object}  response.ErrorResponse
// @Router       /api/v1/admin/users/{id}/plan [patch]
func (h *AdminHandler) UpdateUserPlan(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(chi.URLParam(r, "id"))

	var req UpdateUserPlanReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "bad request")
		return
	}

	if err := h.userService.ChangePlan(r.Context(), id, req.PlanID); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "updated"})
}
