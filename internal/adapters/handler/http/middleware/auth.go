package middleware

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shanth1/gotrace/internal/config"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

func JWTAuth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				response.Error(w, http.StatusUnauthorized, "Invalid token format")
				return
			}

			claims := &domain.JwtCustomClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.Auth.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), domain.CtxKeyUser, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(domain.CtxKeyUser).(*domain.JwtCustomClaims)
		if !ok {
			response.Error(w, http.StatusUnauthorized, "User not authorized")
			return
		}

		if claims.Role != domain.RoleAdmin {
			response.Error(w, http.StatusForbidden, "Access denied: admins only")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func BasicAuth(user, pass string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if !ok {
				basicAuthFailed(w)
				return
			}

			userMatch := subtle.ConstantTimeCompare([]byte(u), []byte(user)) == 1
			passMatch := subtle.ConstantTimeCompare([]byte(p), []byte(pass)) == 1

			if !userMatch || !passMatch {
				basicAuthFailed(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func basicAuthFailed(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	response.Error(w, http.StatusUnauthorized, "unauthorized")
}
