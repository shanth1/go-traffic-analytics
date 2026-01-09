package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestCampaignHandler_GetCampaigns(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCampaignService := mocks.NewMockCampaignService(ctrl)
	handler := NewCampaignHandler(mockCampaignService)

	t.Run("success", func(t *testing.T) {
		userID := domain.UserID("user123")
		expectedCampaigns := []*domain.Campaign{
			{ID: "camp1", UserID: userID, Name: "Campaign 1"},
		}
		expectedCampaignsCount := int64(len(expectedCampaigns))

		limit := 2
		mockCampaignService.EXPECT().
			GetCampaigns(gomock.Any(), domain.CampaignFilter{
				UserID: userID,
				Limit:  limit,
			}).
			Return(expectedCampaigns, expectedCampaignsCount, nil).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns?limit="+fmt.Sprintf("%d", limit), nil)
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetCampaigns(w, req)

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
		limit := 2

		expectedCampaignCount := int64(0)
		mockCampaignService.EXPECT().
			GetCampaigns(gomock.Any(), domain.CampaignFilter{
				UserID: userID,
				Limit:  limit,
			}).
			Return(nil, expectedCampaignCount, errors.New("service error")).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/campaigns?limit="+fmt.Sprintf("%d", limit), nil)
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.GetCampaigns(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestCampaignHandler_CreateCampaign(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCampaignService := mocks.NewMockCampaignService(ctrl)
	handler := NewCampaignHandler(mockCampaignService)

	t.Run("success", func(t *testing.T) {
		userID := domain.UserID("user123")
		reqBody := CreateCampaignRequest{Name: "New Campaign"}
		body, _ := json.Marshal(reqBody)

		expectedCampaign := &domain.Campaign{
			ID:   "camp123",
			Name: "New Campaign",
		}

		mockCampaignService.EXPECT().
			CreateCampaign(gomock.Any(), userID, "New Campaign").
			Return(expectedCampaign, nil).
			Times(1)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", bytes.NewReader(body))
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateCampaign(w, req)

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

	t.Run("bad request - missing name", func(t *testing.T) {
		userID := domain.UserID("user123")
		reqBody := CreateCampaignRequest{}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", bytes.NewReader(body))
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateCampaign(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		userID := domain.UserID("user123")
		reqBody := CreateCampaignRequest{Name: "New Campaign"}
		body, _ := json.Marshal(reqBody)

		mockCampaignService.EXPECT().
			CreateCampaign(gomock.Any(), userID, "New Campaign").
			Return(nil, errors.New("service error")).
			Times(1)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns", bytes.NewReader(body))
		claims := &domain.JwtCustomClaims{UserID: userID}
		ctx := context.WithValue(req.Context(), domain.CtxKeyUser, claims)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler.CreateCampaign(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}
