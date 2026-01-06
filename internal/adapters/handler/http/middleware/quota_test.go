package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestQuotaMiddleware_CheckClickLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLinkRepo := mocks.NewMockLinkRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPlanRepo := mocks.NewMockPlanRepository(ctrl)

	mw := NewQuotaMiddleware(mockLinkRepo, mockUserRepo, mockPlanRepo)

	t.Run("within limit", func(t *testing.T) {
		link := &domain.Link{UserID: "user123"}
		user := &domain.User{ID: "user123", ClicksCurrentMonth: 10}
		plan := &domain.Plan{MaxClicksMonth: 100}

		mockLinkRepo.EXPECT().FindBySlug(gomock.Any(), "abc").Return(link, nil)
		mockUserRepo.EXPECT().FindByID(gomock.Any(), domain.UserID("user123")).Return(user, nil)
		mockPlanRepo.EXPECT().FindByID(gomock.Any(), user.PlanID).Return(plan, nil)

		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			nextCalled = true
		})

		req := httptest.NewRequest(http.MethodGet, "/abc", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("slug", "abc")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		mw.CheckClickLimit(next).ServeHTTP(w, req)

		if !nextCalled {
			t.Errorf("next handler should have been called")
		}
	})

	t.Run("limit exceeded", func(t *testing.T) {
		link := &domain.Link{UserID: "user123"}
		user := &domain.User{ID: "user123", ClicksCurrentMonth: 100}
		plan := &domain.Plan{MaxClicksMonth: 100}

		mockLinkRepo.EXPECT().FindBySlug(gomock.Any(), "abc").Return(link, nil)
		mockUserRepo.EXPECT().FindByID(gomock.Any(), domain.UserID("user123")).Return(user, nil)
		mockPlanRepo.EXPECT().FindByID(gomock.Any(), user.PlanID).Return(plan, nil)

		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			nextCalled = true
		})

		req := httptest.NewRequest(http.MethodGet, "/abc", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("slug", "abc")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		mw.CheckClickLimit(next).ServeHTTP(w, req)

		if !nextCalled {
			t.Errorf("next handler should have been called even if limit exceeded (for now)")
		}
	})
}
