package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports/mocks"
	"go.uber.org/mock/gomock"
)

func TestBillingHandler_GetPlans(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBillingService := mocks.NewMockBillingService(ctrl)
	handler := NewBillingHandler(mockBillingService)

	t.Run("success", func(t *testing.T) {
		expectedPlans := []*domain.Plan{
			{ID: "free", Name: "Free"},
		}

		mockBillingService.EXPECT().
			GetAllPlans(gomock.Any()).
			Return(expectedPlans, nil).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/billing/plans", nil)
		w := httptest.NewRecorder()

		handler.GetPlans(w, req)

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
		mockBillingService.EXPECT().
			GetAllPlans(gomock.Any()).
			Return(nil, errors.New("service error")).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/billing/plans", nil)
		w := httptest.NewRecorder()

		handler.GetPlans(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}
