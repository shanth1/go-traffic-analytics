package memory

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

type Cache struct {
	mu    sync.RWMutex
	items map[string]item
}

func NewCache() ports.Cache {
	return &Cache{
		items: make(map[string]item),
	}
}

func (c *Cache) Set(_ context.Context, key string, value interface{}, ttl time.Duration) error {
	const op = "memory.Cache.Set"

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}

	return nil
}

func (c *Cache) Get(_ context.Context, key string) (interface{}, error) {
	const op = "memory.Cache.Get"

	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		return nil, ops.Wrap(op, ops.KindNotFound, ports.ErrCacheMiss)
	}

	if time.Now().After(item.expiresAt) {
		go c.Delete(context.Background(), key)
		return nil, ops.Wrap(op, ops.KindOther, ports.ErrCacheMiss)
	}

	return item.value, nil
}

func (c *Cache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)

	return nil
}
