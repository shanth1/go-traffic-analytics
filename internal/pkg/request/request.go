package request

import (
	"net/http"

	"github.com/shanth1/gotrace/internal/core/domain"
)

func GetUserClaims(r *http.Request) *domain.JwtCustomClaims {
	claims, ok := r.Context().Value(domain.CtxKeyUser).(*domain.JwtCustomClaims)
	if !ok {
		return nil
	}
	return claims
}

func GetUserID(r *http.Request) string {
	claims := GetUserClaims(r)
	if claims == nil {
		return ""
	}
	return claims.UserID
}
