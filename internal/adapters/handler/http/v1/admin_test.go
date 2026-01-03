package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports/mocks"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
	"go.uber.org/mock/gomock"
)

func TestAdminHandler_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewAdminHandler(mockUserService)

	t.Run("success", func(t *testing.T) {
		expectedUsers := []*domain.User{
			{ID: "user1", Email: "user1@example.com"},
		}

		mockUserService.EXPECT().
			GetAll(gomock.Any(), 1, 20).
			Return(expectedUsers, nil).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users?page=1", nil)
		w := httptest.NewRecorder()

		handler.GetUsers(w, req)

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
		mockUserService.EXPECT().
			GetAll(gomock.Any(), 1, 20).
			Return(nil, errors.New("service error")).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users?page=1", nil)
		w := httptest.NewRecorder()

		handler.GetUsers(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestAdminHandler_UpdateUserStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewAdminHandler(mockUserService)

	t.Run("success", func(t *testing.T) {
		userID := domain.UserID("user123")
		reqBody := UpdateUserStatusReq{IsActive: true}
		body, _ := json.Marshal(reqBody)

		mockUserService.EXPECT().
			SetStatus(gomock.Any(), userID, true).
			Return(nil).
			Times(1)

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/users/user123/status", bytes.NewReader(body))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "user123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.UpdateUserStatus(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}
		if data, ok := resp["data"].(map[string]interface{}); ok {
			if data["status"] != "updated" {
				t.Errorf("expected status 'updated', got %v", data["status"])
			}
		} else {
			t.Errorf("expected data key in response")
		}
	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/users/user123/status", bytes.NewReader([]byte("invalid")))
		w := httptest.NewRecorder()

		handler.UpdateUserStatus(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		userID := domain.UserID("user123")
		reqBody := UpdateUserStatusReq{IsActive: true}
		body, _ := json.Marshal(reqBody)

		mockUserService.EXPECT().
			SetStatus(gomock.Any(), userID, true).
			Return(errors.New("service error")).
			Times(1)

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/users/user123/status", bytes.NewReader(body))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "user123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.UpdateUserStatus(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestAdminHandler_UpdateUserPlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	handler := NewAdminHandler(mockUserService)

	t.Run("success", func(t *testing.T) {
		userID := domain.UserID("user123")
		reqBody := UpdateUserPlanReq{PlanID: "pro"}
		body, _ := json.Marshal(reqBody)

		mockUserService.EXPECT().
			ChangePlan(gomock.Any(), userID, "pro").
			Return(nil).
			Times(1)

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/users/user123/plan", bytes.NewReader(body))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "user123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.UpdateUserPlan(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp response.DataResponse[struct{ Status string }]
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if resp.Data.Status != "updated" {
			t.Errorf("expected status 'updated', got %v", resp.Data.Status)
		}

	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/users/user123/plan", bytes.NewReader([]byte("invalid")))
		w := httptest.NewRecorder()

		handler.UpdateUserPlan(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		userID := domain.UserID("user123")
		reqBody := UpdateUserPlanReq{PlanID: "pro"}
		body, _ := json.Marshal(reqBody)

		mockUserService.EXPECT().
			ChangePlan(gomock.Any(), userID, "pro").
			Return(errors.New("service error")).
			Times(1)

		req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/users/user123/plan", bytes.NewReader(body))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "user123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.UpdateUserPlan(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}
