package v1

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/shanth1/gotrace/internal/core/services"
)

type AdminHandler struct {
	userService *services.UserService // Нужен отдельный сервис для управления юзерами
	planService *services.PlanService // Для чтения планов (если он выделен)
}

func NewAdminHandler(u *services.UserService) *AdminHandler {
	return &AdminHandler{userService: u}
}

func (h *AdminHandler) GetUsers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	limit := 20

	users, err := h.userService.GetAllUsers(c.Request().Context(), page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": users})
}

type updateUserStatusReq struct {
	IsActive bool `json:"is_active"`
}

func (h *AdminHandler) UpdateUserStatus(c echo.Context) error {
	id := c.Param("id")
	var req updateUserStatusReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	if err := h.userService.SetUserStatus(c.Request().Context(), id, req.IsActive); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

type updateUserPlanReq struct {
	PlanID string `json:"plan_id"`
}

func (h *AdminHandler) UpdateUserPlan(c echo.Context) error {
	id := c.Param("id")
	var req updateUserPlanReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	if err := h.userService.ChangeUserPlan(c.Request().Context(), id, req.PlanID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

func (h *AdminHandler) GetPlans(c echo.Context) error {
	// Предполагаем наличие метода в userService или planService
	plans, err := h.userService.GetAllPlans(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": plans})
}
