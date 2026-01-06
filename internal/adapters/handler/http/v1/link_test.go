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
	"go.uber.org/mock/gomock"
)

func TestLinkHandler_GetLinks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLinkService := mocks.NewMockLinkService(ctrl)
	mockCampaignService := mocks.NewMockCampaignService(ctrl)
	handler := NewLinkHandler(mockLinkService, mockCampaignService)

	t.Run("success", func(t *testing.T) {
		userID := domain.UserID("user123")
		expectedLinks := []*domain.Link{
			{ID: "link1", Slug: "abc", TargetURL: "http://example.com"},
		}

		mockLinkService.EXPECT().
			GetLinkList(gomock.Any(), gomock.Any()).
			Return(expectedLinks, nil).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/links", nil)
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetLinks(w, req)

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

		mockLinkService.EXPECT().
			GetLinkList(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("service error")).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/links", nil)
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetLinks(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestLinkHandler_CreateLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLinkService := mocks.NewMockLinkService(ctrl)
	mockCampaignService := mocks.NewMockCampaignService(ctrl)
	handler := NewLinkHandler(mockLinkService, mockCampaignService)

	t.Run("success", func(t *testing.T) {
		userID := domain.UserID("user123")
		reqBody := CreateLinkRequest{
			CampaignID: "camp1",
			TargetURL:  "http://example.com",
		}
		body, _ := json.Marshal(reqBody)

		expectedLink := &domain.Link{
			ID:        "link123",
			Slug:      "abc",
			TargetURL: "http://example.com",
		}

		mockLinkService.EXPECT().
			CreateLink(gomock.Any(), gomock.Any()).
			Return(expectedLink, nil).
			Times(1)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewReader(body))
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateLink(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}
		if resp["data"] == nil {
			t.Errorf("expected data in response")
		}
	})

	t.Run("bad request - missing target_url", func(t *testing.T) {
		userID := domain.UserID("user123")
		reqBody := CreateLinkRequest{CampaignID: "camp1"}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewReader(body))
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateLink(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		userID := domain.UserID("user123")
		reqBody := CreateLinkRequest{
			CampaignID: "camp1",
			TargetURL:  "http://example.com",
		}
		body, _ := json.Marshal(reqBody)

		mockLinkService.EXPECT().
			CreateLink(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("service error")).
			Times(1)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewReader(body))
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateLink(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestLinkHandler_DeleteLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLinkService := mocks.NewMockLinkService(ctrl)
	mockCampaignService := mocks.NewMockCampaignService(ctrl)
	handler := NewLinkHandler(mockLinkService, mockCampaignService)

	t.Run("success", func(t *testing.T) {
		userID := domain.UserID("user123")
		linkID := domain.LinkID("link123")

		mockLinkService.EXPECT().
			DeleteLink(gomock.Any(), userID, linkID).
			Return(nil).
			Times(1)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/links/link123", nil)
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "link123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.DeleteLink(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		userID := domain.UserID("user123")
		linkID := domain.LinkID("link123")

		mockLinkService.EXPECT().
			DeleteLink(gomock.Any(), userID, linkID).
			Return(errors.New("service error")).
			Times(1)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/links/link123", nil)
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "link123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		handler.DeleteLink(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}
