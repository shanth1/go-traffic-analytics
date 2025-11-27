package domain

import "github.com/golang-jwt/jwt/v5"

type ContextKey string

const (
	CtxKeyUser ContextKey = "user_claims"
)

type JwtCustomClaims struct {
	UserID string   `json:"id"`
	Email  string   `json:"email"`
	Role   UserRole `json:"role"`
	jwt.RegisteredClaims
}
