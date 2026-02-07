package cachedproxy

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotools/logkeys"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type LinkRepo struct {
	repo   ports.LinkRepository
	cache  ports.Cache
	logger log.Logger
	ttl    time.Duration
}

func NewLinkRepo(repo ports.LinkRepository, cache ports.Cache, logger log.Logger, ttl time.Duration) ports.LinkRepository {
	return &LinkRepo{
		repo:   repo,
		cache:  cache,
		logger: logger,
		ttl:    ttl,
	}
}

func (r *LinkRepo) buildKey(criteria string, value string) string {
	return fmt.Sprintf("gotrace:link:%s:%s", criteria, value)
}

func (r *LinkRepo) Save(ctx context.Context, link *domain.Link) error {
	if err := r.repo.Save(ctx, link); err != nil {
		return err
	}

	r.invalidate(string(link.ID), link.Slug)

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
	return r.repo.FindAll(ctx, filter)
}

func (r *LinkRepo) Count(ctx context.Context, filter domain.LinkFilter) (int64, error) {
	return r.repo.Count(ctx, filter)
}

func (r *LinkRepo) Delete(ctx context.Context, userID domain.UserID, id domain.LinkID) error {
	var slugToDelete string

	keyID := r.buildKey("id", string(id))
	cachedLink, _ := r.getFromCache(ctx, keyID)
	if cachedLink != nil {
		slugToDelete = cachedLink.Slug
	} else {
		if l, err := r.repo.FindByID(ctx, id); err == nil {
			slugToDelete = l.Slug
		}
	}

	if err := r.repo.Delete(ctx, userID, id); err != nil {
		return err
	}

	r.invalidate(string(id), slugToDelete)

	return nil
}

func (r *LinkRepo) DeleteByUserID(ctx context.Context, userID domain.UserID) error {
	return r.repo.DeleteByUserID(ctx, userID)
}

func (r *LinkRepo) DeleteByCampaignID(ctx context.Context, campaignID string) error {
	return r.repo.DeleteByCampaignID(ctx, campaignID)
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
		r.logger.Error().Err(err).Str(logkeys.CacheKey, key).Msg("unmarshal error in link repo")
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepo) setCache(key string, link *domain.Link) {
	go func() {
		bytes, err := json.Marshal(link)
		if err != nil {
			r.logger.Error().Err(err).Msg("marshal error")
			return
		}

		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.cache.Set(bgCtx, key, bytes, r.ttl); err != nil {
			r.logger.Error().Err(err).Str(logkeys.CacheKey, key).Msg("failed to set cache")
		}
	}()
}

func (r *LinkRepo) invalidate(idStr, slug string) {
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.cache.Delete(bgCtx, r.buildKey("id", idStr)); err != nil {
			r.logger.Error().Err(err).Str("id", idStr).Msg("failed to invalidate link id")
		}

		if slug != "" {
			if err := r.cache.Delete(bgCtx, r.buildKey("slug", slug)); err != nil {
				r.logger.Error().Err(err).Str("slug", slug).Msg("failed to invalidate link slug")
			}
		}
	}()
}
