package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_GetProfileTree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewUserHandler(mockUserService)

	t.Run("success", func(t *testing.T) {
		userID := domain.UserID("user123")
		expectedTree := &domain.HierarchyNode{
			Name: "User",
			Type: "root",
			Children: []*domain.HierarchyNode{
				{Name: "Campaign1", Type: "campaign"},
			},
		}

		mockUserService.EXPECT().
			GetHierarchy(gomock.Any(), userID).
			Return(expectedTree, nil).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/tree", nil)
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetProfileTree(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}
		if resp["data"] == nil {
			t.Errorf("expected data in response")
		}
	})

	t.Run("service error", func(t *testing.T) {
		userID := domain.UserID("user123")

		mockUserService.EXPECT().
			GetHierarchy(gomock.Any(), userID).
			Return(nil, errors.New("service error")).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/tree", nil)
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetProfileTree(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}
