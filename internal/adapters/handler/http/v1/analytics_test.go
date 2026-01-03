package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports/mocks"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
	"go.uber.org/mock/gomock"
)

func TestAnalyticsHandler_GetSummary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockAnalyticsService(ctrl)
	handler := NewAnalyticsHandler(mockService)

	t.Run("success", func(t *testing.T) {
		expectedSummary := &domain.Summary{
			TotalClicks: 100,
			// Fill other fields as needed
		}

		mockService.EXPECT().
			GetSummary(gomock.Any(), gomock.Any()).
			Return(expectedSummary, nil).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/summary", nil)
		w := httptest.NewRecorder()

		handler.GetSummary(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp response.DataResponse[domain.Summary]
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		// TODO:
		if resp.Data.TotalClicks != expectedSummary.TotalClicks {
			t.Error("mismatch data")
		}
	})

	t.Run("service error", func(t *testing.T) {
		mockService.EXPECT().
			GetSummary(gomock.Any(), gomock.Any()).
			Return(nil, errors.New("service error")).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/summary", nil)
		w := httptest.NewRecorder()

		handler.GetSummary(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})

	// TODO:
	t.Run("with filters", func(t *testing.T) {
		from := time.Now().AddDate(0, 0, -7).Round(time.Minute)
		to := time.Now().Round(time.Minute)

		expectedSummary := &domain.Summary{
			TotalClicks: 50,
		}

		mockService.EXPECT().
			GetSummary(gomock.Any(), gomock.Any()).
			Return(expectedSummary, nil).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/analytics/summary?campaign_id=test_campaign&link_id=test_link&from="+from.Format(time.RFC3339)+"&to="+to.Format(time.RFC3339), nil)
		w := httptest.NewRecorder()

		handler.GetSummary(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}
