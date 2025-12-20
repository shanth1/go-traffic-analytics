package ports

import (
	"context"
	"errors"
	"time"
)

//go:generate mockgen -source=cache.go -destination=mocks/cache_mock.go -package=mocks

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

var ErrCacheMiss = errors.New("cache miss")
