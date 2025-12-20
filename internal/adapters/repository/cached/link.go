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

func (r *LinkRepo) buildKey(criteria, value string) string {
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
			r.buildKey("id", link.ID),
			r.buildKey("slug", link.Slug),
		}

		for _, key := range keys {
			// TODO: logging
			_ = r.cache.Delete(bgCtx, key)
		}
	}()

	return nil
}

func (r *LinkRepo) FindBySlug(ctx context.Context, slug string) (*domain.Link, error) {
	key := r.buildKey("slug", slug)

	val, err := r.cache.Get(ctx, key)
	if err == nil {
		if bytesVal, ok := val.([]byte); ok {
			var link domain.Link
			if jsonErr := json.Unmarshal(bytesVal, &link); jsonErr == nil {
				return &link, nil
			}
			// TODO: logging
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

		// TODO: logging
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return link, nil
}

func (r *LinkRepo) FindAll(ctx context.Context) ([]*domain.Link, error) {
	links, err := r.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (r *LinkRepo) FindAllByCampaignID(ctx context.Context, campaignID string) ([]*domain.Link, error) {
	links, err := r.repo.FindAllByCampaignID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func (r *LinkRepo) CountByUserID(ctx context.Context, userID string) (int64, error) {
	count, err := r.repo.CountByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
