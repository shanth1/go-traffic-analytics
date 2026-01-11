package middleware

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotools/logkeys"
	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/config"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

func JWTAuth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.JWTAuth"

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				err := ops.New(op, ops.KindUnauthorized, "missing authorization header")
				response.Error(w, r, err)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				err := ops.New(op, ops.KindUnauthorized, "invalid token format")
				response.Error(w, r, err)
				return
			}

			claims := &domain.JwtCustomClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.Auth.JWTSecret), nil
			})
			if err != nil {
				err := ops.WrapMsg(op, ops.KindUnauthorized, err, "invalid or expired token")
				response.Error(w, r, err)
				return
			}
			if !token.Valid {
				response.Error(w, r, ops.New(op, ops.KindUnauthorized, "token is invalid"))
				response.Error(w, r, err)
				return
			}

			userID := claims.UserID

			logger := log.FromContext(r.Context())
			enrichedLogger := logger.With(log.Str(logkeys.UserID, string(userID)))
			ctx := log.NewContext(r.Context(), enrichedLogger)
			ctx = context.WithValue(ctx, domain.CtxKeyUser, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "middleware.AdminOnly"

		claims, ok := r.Context().Value(domain.CtxKeyUser).(*domain.JwtCustomClaims)
		if !ok {
			err := ops.New(op, ops.KindUnauthorized, "user not authorized")
			response.Error(w, r, err)
			return
		}

		if claims.Role != domain.RoleAdmin {
			err := ops.New(op, ops.KindPermission, "access denied")
			response.Error(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func APIKeyAuth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.AdminOnly"

			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				err := ops.New(op, ops.KindUnauthorized, "missing authorization header")
				response.Error(w, r, err)
				return
			}

			if subtle.ConstantTimeCompare([]byte(apiKey), []byte(cfg.Auth.APIKey)) != 1 {
				err := ops.New(op, ops.KindUnauthorized, "invalid api token")
				response.Error(w, r, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func BasicAuth(user, pass string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if !ok {
				basicAuthFailed(w, r)
				return
			}

			userMatch := subtle.ConstantTimeCompare([]byte(u), []byte(user)) == 1
			passMatch := subtle.ConstantTimeCompare([]byte(p), []byte(pass)) == 1

			if !userMatch || !passMatch {
				basicAuthFailed(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func basicAuthFailed(w http.ResponseWriter, r *http.Request) {
	const op = "middleware.basicAuthFailed"

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	err := ops.New(op, ops.KindUnauthorized, "unauthorized")
	response.Error(w, r, err)
}
