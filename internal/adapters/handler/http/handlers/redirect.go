package handlers

import (
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

type RedirectHandler struct {
	service ports.RedirectService
}

func NewRedirectHandler(s ports.RedirectService) *RedirectHandler {
	return &RedirectHandler{service: s}
}

// Redirect godoc
// @Summary      Redirect link
// @Description  Redirect user to target URL based on slug
// @Tags         Redirect
// @Param        slug  path  string  true  "Short Link Slug"
// @Success      307   {string}  string  "Temporary Redirect"
// @Failure      404   {object}  response.ErrorResponse "Link not found"
// @Router       /{slug} [get]
func (h *RedirectHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	ip := getRealIP(r)
	ua := r.UserAgent()
	referer := r.Referer()

	targetURL, err := h.service.Process(r.Context(), slug, domain.RequestMetadata{
		IP:        ip,
		UserAgent: ua,
		Referer:   referer,
	})
	if err != nil {
		response.Error(w, http.StatusNotFound, "Link not found or inactive")
		return
	}

	http.Redirect(w, r, targetURL, http.StatusTemporaryRedirect)
}

// --- Helpers ---

func getRealIP(r *http.Request) string {
	// Проверяем X-Forwarded-For
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0]) // (client, proxy1, proxy2)
	}

	// X-Real-IP
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}

	// Fallback на RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
