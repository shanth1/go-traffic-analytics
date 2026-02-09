package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shanth1/gotrace/internal/config"
	"github.com/shanth1/gotrace/internal/core/domain"
)

func TestJWTAuth(t *testing.T) {
	cfg := &config.Config{
		Auth: config.Auth{
			JWTSecret: "secret",
		},
	}

	mw := JWTAuth(cfg)

	t.Run("valid token", func(t *testing.T) {
		claims := domain.JwtCustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenStr, err := token.SignedString([]byte("secret"))
		if err != nil {
			t.Errorf("signed string: %v", err)
		}

		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			nextCalled = true
			_, ok := r.Context().Value(domain.CtxKeyUser).(*domain.JwtCustomClaims)
			if !ok {
				t.Errorf("expected user claims in context")
			}
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenStr)
		w := httptest.NewRecorder()

		mw(next).ServeHTTP(w, req)

		if !nextCalled {
			t.Errorf("next handler should have been called")
		}
	})

	t.Run("missing api key", func(t *testing.T) {
		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			nextCalled = true
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		mw(next).ServeHTTP(w, req)

		if nextCalled {
			t.Errorf("next handler should not have been called")
		}
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestAdminOnly(t *testing.T) {
	mw := AdminOnly

	t.Run("admin user", func(t *testing.T) {
		claims := &domain.JwtCustomClaims{Role: domain.RoleAdmin}
		ctx := context.WithValue(context.Background(), domain.CtxKeyUser, claims)

		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			nextCalled = true
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		mw(next).ServeHTTP(w, req)

		if !nextCalled {
			t.Errorf("next handler should have been called")
		}
	})

	t.Run("non-admin user", func(t *testing.T) {
		claims := &domain.JwtCustomClaims{Role: domain.RoleClient}
		ctx := context.WithValue(context.Background(), domain.CtxKeyUser, claims)

		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			nextCalled = true
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		w := httptest.NewRecorder()

		mw(next).ServeHTTP(w, req)

		if nextCalled {
			t.Errorf("next handler should not have been called")
		}
		if w.Code != http.StatusForbidden {
			t.Errorf("expected status %d, got %d", http.StatusForbidden, w.Code)
		}
	})
}

func TestAPIKeyAuth(t *testing.T) {
	cfg := &config.Config{
		Auth: config.Auth{
			APIKey: "valid-api-key",
		},
	}

	mw := APIKeyAuth(cfg)

	t.Run("valid api key", func(t *testing.T) {
		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			nextCalled = true
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-API-Key", "valid-api-key")
		w := httptest.NewRecorder()

		mw(next).ServeHTTP(w, req)

		if !nextCalled {
			t.Errorf("next handler should have been called")
		}
	})

	t.Run("invalid api key", func(t *testing.T) {
		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			nextCalled = true
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-API-Key", "invalid")
		w := httptest.NewRecorder()

		mw(next).ServeHTTP(w, req)

		if nextCalled {
			t.Errorf("next handler should not have been called")
		}
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestBasicAuth(t *testing.T) {
	mw := BasicAuth("user", "pass")

	t.Run("valid credentials", func(t *testing.T) {
		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			nextCalled = true
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetBasicAuth("user", "pass")
		w := httptest.NewRecorder()

		mw(next).ServeHTTP(w, req)

		if !nextCalled {
			t.Errorf("next handler should have been called")
		}
	})

	t.Run("invalid credentials", func(t *testing.T) {
		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			nextCalled = true
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetBasicAuth("user", "wrong")
		w := httptest.NewRecorder()

		mw(next).ServeHTTP(w, req)

		if nextCalled {
			t.Errorf("next handler should not have been called")
		}
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestJWTAuthOptional(t *testing.T) {
	cfg := &config.Config{
		Auth: config.Auth{
			JWTSecret: "secret",
		},
	}

	mw := JWTAuthOptional(cfg)

	t.Run("no auth header (anonymous)", func(t *testing.T) {
		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			nextCalled = true
			// Проверяем, что в контексте НЕТ юзера
			_, ok := r.Context().Value(domain.CtxKeyUser).(*domain.JwtCustomClaims)
			if ok {
				t.Error("expected NO user claims in context for anonymous request")
			}
		})

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()

		mw(next).ServeHTTP(w, req)

		if !nextCalled {
			t.Errorf("next handler should have been called")
		}
		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("valid token (authenticated)", func(t *testing.T) {
		userID := domain.UserID("user-123")
		claims := domain.JwtCustomClaims{
			UserID: userID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, _ := token.SignedString([]byte(cfg.Auth.JWTSecret))

		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			nextCalled = true
			extractedClaims, ok := r.Context().Value(domain.CtxKeyUser).(*domain.JwtCustomClaims)
			if !ok {
				t.Error("expected user claims in context")
				return
			}
			if extractedClaims.UserID != userID {
				t.Errorf("expected user ID %s, got %s", userID, extractedClaims.UserID)
			}
		})

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenStr)
		w := httptest.NewRecorder()

		mw(next).ServeHTTP(w, req)

		if !nextCalled {
			t.Errorf("next handler should have been called")
		}
	})

	t.Run("invalid token (fallback to anonymous)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-garbage-token")
		w := httptest.NewRecorder()

		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			nextCalled = true
			// В контексте не должно быть юзера, так как токен битый
			_, ok := r.Context().Value(domain.CtxKeyUser).(*domain.JwtCustomClaims)
			if ok {
				t.Error("expected NO user claims for invalid token")
			}
		})

		mw(next).ServeHTTP(w, req)

		if !nextCalled {
			t.Errorf("next handler should have been called even with invalid token")
		}
		if w.Code != http.StatusOK {
			t.Errorf("expected status %d (pass-through), got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("expired token (fallback to anonymous)", func(t *testing.T) {
		claims := domain.JwtCustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Протух час назад
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, _ := token.SignedString([]byte(cfg.Auth.JWTSecret))

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenStr)
		w := httptest.NewRecorder()

		nextCalled := false
		next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			nextCalled = true
			_, ok := r.Context().Value(domain.CtxKeyUser).(*domain.JwtCustomClaims)
			if ok {
				t.Error("expected NO user claims for expired token")
			}
		})

		mw(next).ServeHTTP(w, req)

		if !nextCalled {
			t.Errorf("next handler should have been called")
		}
	})
}
