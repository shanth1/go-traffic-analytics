package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	HealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp response.ResponseWrapper[struct {
		Status string `json:"status"`
	}]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Data.Status != "OK" {
		t.Errorf("expected status %s, got %s", "OK", resp.Data.Status)
	}
}
