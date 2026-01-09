package cachedproxy

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type LinkRepo struct {
	repo  ports.LinkRepository
	cache ports.Cache
	ttl   time.Duration
}

func NewLinkRepo(repo ports.LinkRepository, cache ports.Cache, ttl time.Duration) ports.LinkRepository {
	return &LinkRepo{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (r *LinkRepo) buildKey(criteria string, value string) string {
	return fmt.Sprintf("gotrace:link:%s:%s", criteria, value)
}

func (r *LinkRepo) Save(ctx context.Context, link *domain.Link) error {
	if err := r.repo.Save(ctx, link); err != nil {
		return err
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		keys := []string{
			r.buildKey("id", string(link.ID)),
			r.buildKey("slug", link.Slug),
		}

		for _, key := range keys {
			// TODO: logging
			_ = r.cache.Delete(bgCtx, key)
		}
	}()

	return nil
}

func (r *LinkRepo) FindByID(ctx context.Context, id domain.LinkID) (*domain.Link, error) {
	key := r.buildKey("id", string(id))

	if link, err := r.getFromCache(ctx, key); err == nil && link != nil {
		return link, nil
	}

	link, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	r.setCache(key, link)

	return link, nil
}

func (r *LinkRepo) FindBySlug(ctx context.Context, slug string) (*domain.Link, error) {
	key := r.buildKey("slug", slug)
	if link, err := r.getFromCache(ctx, key); err == nil && link != nil {
		return link, nil
	}
	link, err := r.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	r.setCache(key, link)
	return link, nil
}

func (r *LinkRepo) FindAll(ctx context.Context, filter domain.LinkFilter) ([]*domain.Link, error) {
	links, err := r.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (r *LinkRepo) Count(ctx context.Context, filter domain.LinkFilter) (int64, error) {
	return r.repo.Count(ctx, filter)
}

func (r *LinkRepo) Delete(ctx context.Context, userID domain.UserID, id domain.LinkID) error {
	link, err := r.repo.FindByID(ctx, id)

	var slugToDelete string
	if err == nil && link != nil {
		slugToDelete = link.Slug
	}

	if err := r.repo.Delete(ctx, userID, id); err != nil {
		return err
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.cache.Delete(bgCtx, r.buildKey("id", string(id)))
		if slugToDelete != "" {
			_ = r.cache.Delete(bgCtx, r.buildKey("slug", slugToDelete))
		}
	}()

	return nil
}

func (r *LinkRepo) getFromCache(ctx context.Context, key string) (*domain.Link, error) {
	val, err := r.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	bytesVal, ok := val.([]byte)
	if !ok {
		return nil, nil
	}
	var link domain.Link
	if err := json.Unmarshal(bytesVal, &link); err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepo) setCache(key string, link *domain.Link) {
	go func() {
		bytes, err := json.Marshal(link)
		if err == nil {
			bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
		}
	}()
}
