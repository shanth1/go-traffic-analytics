package cachememory

import (
	"context"
	"sync"
	"time"

	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type item struct {
	value     interface{}
	expiresAt time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]item
}

func NewCache() ports.Cache {
	return &MemoryCache{
		items: make(map[string]item),
	}
}

func (c *MemoryCache) Set(_ context.Context, key string, value interface{}, ttl time.Duration) error {
	const op = "cachememory.MemoryCache.Set"

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}

	return nil
}

func (c *MemoryCache) Get(_ context.Context, key string) (interface{}, error) {
	const op = "cachememory.MemoryCache.Get"

	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, ops.E(op, ports.ErrCacheMiss)
	}

	if time.Now().After(item.expiresAt) {
		delete(c.items, key)
		return nil, ops.E(op, ports.ErrCacheMiss)
	}

	return item.value, nil
}

func (c *MemoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)

	return nil
}
