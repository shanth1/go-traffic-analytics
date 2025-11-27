package request

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

func getUserIDFromContext(ctx context.Context) string {
	token, ok := ctx.Value(UserContextKey).(*jwt.Token)
	if !ok {
		return ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ""
	}

	if id, ok := claims["id"].(string); ok {
		return id
	}

	return ""
}

func GetUserID(r *http.Request) string {
	return getUserIDFromContext(r.Context())
}
