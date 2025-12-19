package memory

import (
	"context"
	"sync"
	"time"

	"github.com/shanth1/gotrace/internal/core/ports"
)

type item struct {
	value     interface{}
	expiresAt time.Time
}

type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]item
}

func NewCache() ports.Cache {
	return &InMemoryCache{
		items: make(map[string]item),
	}
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

func (c *InMemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, ports.ErrCacheMiss
	}

	// Проверяем протухание (TTL)
	if time.Now().After(item.expiresAt) {
		// В реальном Redis это делается автоматически,
		// здесь мы лениво удаляем при чтении или можно запустить горутину очистки
		delete(c.items, key)
		return nil, ports.ErrCacheMiss
	}

	return item.value, nil
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}
