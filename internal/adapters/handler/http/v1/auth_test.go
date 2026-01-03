package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports/mocks"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
	"go.uber.org/mock/gomock"
)

func TestAuthHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService)

	t.Run("success", func(t *testing.T) {
		reqBody := RegisterReq{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		expectedUser := &domain.User{
			ID:    "user123",
			Email: "test@example.com",
			Role:  domain.RoleClient,
		}

		mockAuthService.EXPECT().
			Register(gomock.Any(), "test@example.com", "password123").
			Return(expectedUser, nil).
			Times(1)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("X-API-Key", "some-api-key") // Assuming API key auth
		w := httptest.NewRecorder()

		handler.Register(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var resp response.DataResponse[domain.User]
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if resp.Data.ID != "user123" {
			t.Errorf("expected user ID 'user123', got %v", resp.Data.ID)
		}

	})

	t.Run("service error", func(t *testing.T) {
		reqBody := RegisterReq{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		mockAuthService.EXPECT().
			Register(gomock.Any(), "test@example.com", "password123").
			Return(nil, errors.New("registration failed")).
			Times(1)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("X-API-Key", "some-api-key")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
		}
	})

	t.Run("bad request - missing email", func(t *testing.T) {
		reqBody := RegisterReq{
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("X-API-Key", "some-api-key")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthService(ctrl)
	handler := NewAuthHandler(mockAuthService)

	t.Run("success", func(t *testing.T) {
		reqBody := LoginReq{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		expectedToken := "jwt-token"
		expectedUser := &domain.User{
			ID:    "user123",
			Email: "test@example.com",
		}

		mockAuthService.EXPECT().
			Login(gomock.Any(), "test@example.com", "password123").
			Return(expectedToken, expectedUser, nil).
			Times(1)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Login(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp response.DataResponse[struct {
			Token string
			User  domain.User
		}]
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if resp.Data.Token != expectedToken {
			t.Errorf("expected token %s, got %v", expectedToken, resp.Data.Token)
		}
		if resp.Data.User.ID != "user123" {
			t.Errorf("expected user ID 'user123', got %v", resp.Data.User.ID)
		}

	})

	t.Run("service error", func(t *testing.T) {
		reqBody := LoginReq{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		mockAuthService.EXPECT().
			Login(gomock.Any(), "test@example.com", "password123").
			Return("", nil, errors.New("invalid credentials")).
			Times(1)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Login(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.Login(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
