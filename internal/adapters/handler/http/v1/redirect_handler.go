package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shanth1/gotrace/internal/core/services"
)

type RedirectHandler struct {
	service *services.RedirectService
}

func NewRedirectHandler(s *services.RedirectService) *RedirectHandler {
	return &RedirectHandler{service: s}
}

func (h *RedirectHandler) Redirect(c echo.Context) error {
	slug := c.Param("slug")

	// Получаем IP и User-Agent для аналитики
	ip := c.RealIP()
	ua := c.Request().UserAgent()
	referer := c.Request().Referer()

	targetURL, err := h.service.ProcessRedirect(c.Request().Context(), slug, ip, ua, referer)
	if err != nil {
		// Можно вернуть красивую HTML страницу 404
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Link not found or inactive"})
	}

	// 307 Temporary Redirect (сохраняет метод POST, если был) или 302 Found
	return c.Redirect(http.StatusTemporaryRedirect, targetURL)
}
