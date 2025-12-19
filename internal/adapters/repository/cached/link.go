package cached

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type CachedLinkRepo struct {
	repo  ports.LinkRepository
	cache ports.Cache
	ttl   time.Duration
}

func NewCachedLinkRepo(repo ports.LinkRepository, cache ports.Cache, ttl time.Duration) ports.LinkRepository {
	return &CachedLinkRepo{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (r *CachedLinkRepo) buildKey(criteria, value string) string {
	// TODO: gotrace -> common consts
	return fmt.Sprintf("gotrace:link:%s:%s", criteria, value)
}

func (r *CachedLinkRepo) Save(ctx context.Context, link *domain.Link) error {
	if err := r.repo.Save(ctx, link); err != nil {
		return err
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		keys := []string{
			r.buildKey("id", link.ID),
			r.buildKey("slug", link.Slug),
		}

		for _, key := range keys {
			if err := r.cache.Delete(bgCtx, key); err != nil {
				// TODO: logging
			}
		}
	}()

	return nil
}

func (r *CachedLinkRepo) FindBySlug(ctx context.Context, slug string) (*domain.Link, error) {
	key := r.buildKey("slug", slug)

	val, err := r.cache.Get(ctx, key)
	if err == nil {
		if bytesVal, ok := val.([]byte); ok {
			var link domain.Link
			if jsonErr := json.Unmarshal(bytesVal, &link); jsonErr == nil {
				return &link, nil
			} else {
				// TODO: logging
			}
		}
	} else if !errors.Is(err, ports.ErrCacheMiss) {
		// TODO: logging
	}

	link, err := r.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	go func() {
		bytes, err := json.Marshal(link)
		if err != nil {
			// TODO: logging
			return
		}

		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.cache.Set(bgCtx, key, bytes, r.ttl); err != nil {
			// TODO: logging
		}
	}()

	return link, nil
}

func (r *CachedLinkRepo) FindAll(ctx context.Context) ([]*domain.Link, error) {
	links, err := r.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (r *CachedLinkRepo) FindAllByCampaignID(ctx context.Context, campaignID string) ([]*domain.Link, error) {
	links, err := r.repo.FindAllByCampaignID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (r *CachedLinkRepo) CountByUserID(ctx context.Context, userID string) (int64, error) {
	count, err := r.repo.CountByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
