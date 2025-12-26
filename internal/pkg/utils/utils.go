package utils

import (
	"context"
	"errors"

	"github.com/shanth1/gotrace/internal/core/domain"
)

func GetUserFromContext(ctx context.Context) (*domain.JwtCustomClaims, error) {
	claims, ok := ctx.Value(domain.CtxKeyUser).(*domain.JwtCustomClaims)
	if !ok || claims == nil {
		return nil, errors.New("no user in context")
	}
	return claims, nil
}
