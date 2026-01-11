package v1

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotools/logkeys"
	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/request"
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

// GetUsers godoc
// @Summary List users
// @Description Get paginated users list (Admin). Returns array wrapped in data with metadata.
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} UsersListResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/admin/users [get]
func (h *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 20
	}

	users, total, err := h.userService.GetAll(r.Context(), page, limit)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	if users == nil {
		users = []*domain.User{}
	}

	publicUsers := make([]domain.UserPublic, len(users))
	for i, user := range users {
		publicUsers[i] = user.ToPublic()
	}

	log.FromContext(r.Context()).Info().Int("page", page).Int("limit", limit).Int64("total", total).Msg("admin_users_listed")

	response.JSON(w, r, http.StatusOK, UsersListResponse{
		Data: publicUsers,
		Meta: domain.PaginationMeta{
			Total:  total,
			Limit:  limit,
			Offset: (page - 1) * limit,
		},
	})
}

// UpdateUserStatus godoc
// @Summary Update user status
// @Description Activate/Deactivate user
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateUserStatusReq true "Status Request"
// @Success 200 {object} StatusResponse "Status: updated"
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/admin/users/{id}/status [patch]
func (h *AdminHandler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	const op = "v1.AdminHandler.UpdateUserStatus"

	id := domain.UserID(chi.URLParam(r, "id"))

	var req UpdateUserStatusReq
	if err := request.DecodeJSON(w, r, &req); err != nil {
		err := ops.WrapMsg(op, ops.KindInvalid, err, err.Error())
		response.Error(w, r, err)
		return
	}

	if err := h.userService.SetStatus(r.Context(), id, req.IsActive); err != nil {
		response.Error(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().Str(logkeys.UserID, string(id)).Bool("is_active", req.IsActive).Msg("user_status_updated")

	response.OK(w, r, StatusData{Status: "updated"})
}

// UpdateUserPlan godoc
// @Summary Change user plan
// @Description Set new plan for user
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateUserPlanReq true "Plan Request"
// @Success 200 {object} StatusResponse "Status: updated"
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/admin/users/{id}/plan [patch]
func (h *AdminHandler) UpdateUserPlan(w http.ResponseWriter, r *http.Request) {
	const op = "v1.AdminHandler.UpdateUserPlan"

	id := domain.UserID(chi.URLParam(r, "id"))

	var req UpdateUserPlanReq
	if err := request.DecodeJSON(w, r, &req); err != nil {
		err := ops.WrapMsg(op, ops.KindInvalid, err, err.Error())
		response.Error(w, r, err)
		return
	}

	if err := h.userService.ChangePlan(r.Context(), id, req.PlanID); err != nil {
		response.Error(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().Str(logkeys.UserID, string(id)).Str("plan_id", req.PlanID).Msg("user_plan_updated")

	response.OK(w, r, StatusData{Status: "updated"})
}
