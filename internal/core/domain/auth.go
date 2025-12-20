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

const (
	ErrUserNotFound = "user not found"
	ErrInvalidCreds = "invalid credentials"
	ErrEmailTaken   = "email already taken"
	ErrInternal     = "internal server error"
	ErrForbidden    = "forbidden"
	ErrUnauthorized = "unauthorized"
	ErrLinkNotFound = "link not found"
	ErrLimitReached = "limit reached"
)
