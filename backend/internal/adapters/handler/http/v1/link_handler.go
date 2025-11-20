package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shanth1/gotrace/internal/core/services"
)

type LinkHandler struct {
	service *services.LinkService
}

func NewLinkHandler(s *services.LinkService) *LinkHandler {
	return &LinkHandler{service: s}
}

// --- Campaigns ---

func (h *LinkHandler) GetCampaigns(c echo.Context) error {
	userID := utils.GetUserIDFromContext(c) // Вспомогательная функция
	campaigns, err := h.service.GetUserCampaigns(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": campaigns})
}

type createCampaignReq struct {
	Name string `json:"name"`
}

func (h *LinkHandler) CreateCampaign(c echo.Context) error {
	userID := utils.GetUserIDFromContext(c)
	var req createCampaignReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	camp, err := h.service.CreateCampaign(c.Request().Context(), userID, req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{"data": camp})
}

// --- Links ---

func (h *LinkHandler) GetLinksByCampaign(c echo.Context) error {
	campaignID := c.Param("id")
	// В реальном сервисе нужно проверить, принадлежит ли кампания этому юзеру!
	links, err := h.service.GetLinks(c.Request().Context(), campaignID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"data": links})
}

type createLinkReq struct {
	CampaignID string `json:"campaign_id"`
	TargetURL  string `json:"target_url"`
	CustomSlug string `json:"custom_slug"` // Optional
}

func (h *LinkHandler) CreateLink(c echo.Context) error {
	userID := utils.GetUserIDFromContext(c)
	var req createLinkReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	link, err := h.service.CreateLink(c.Request().Context(), userID, req.CampaignID, req.TargetURL, req.CustomSlug)
	if err != nil {
		// Если ошибка лимитов, возвращаем 403
		if err == services.ErrLimitReached {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"data": link})
}

func (h *LinkHandler) DeleteLink(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.DeleteLink(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
