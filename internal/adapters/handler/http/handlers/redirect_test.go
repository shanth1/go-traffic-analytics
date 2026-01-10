package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestRedirectHandler_Redirect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedirectService := mocks.NewMockRedirectService(ctrl)
	handler := NewRedirectHandler(mockRedirectService)

	slug := "abc"

	t.Run("success", func(t *testing.T) {
		mockRedirectService.EXPECT().
			Process(gomock.Any(), slug, gomock.Any()).
			Return("http://example.com", nil).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", slug), nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("slug", slug)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.Redirect(w, req)

		if w.Code != http.StatusTemporaryRedirect {
			t.Errorf("expected status %d, got %d", http.StatusTemporaryRedirect, w.Code)
		}

		location := w.Header().Get("Location")
		if location != "http://example.com" {
			t.Errorf("expected location 'http://example.com', got %s", location)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRedirectService.EXPECT().
			Process(gomock.Any(), slug, gomock.Any()).
			Return("", ops.Wrap("test", ops.KindNotFound, errors.New("not found"))).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", slug), nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("slug", slug)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.Redirect(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
