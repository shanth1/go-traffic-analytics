package memoryrepo

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"

	"github.com/shanth1/gotools/errs"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type InMemoryLinkRepo struct {
	mu    sync.RWMutex
	links map[domain.LinkID]*domain.Link
	slugs map[string]domain.LinkID // slug -> id (index)
}

func NewLinkRepo() ports.LinkRepository {
	return &InMemoryLinkRepo{
		links: make(map[domain.LinkID]*domain.Link),
		slugs: make(map[string]domain.LinkID),
	}
}

func (r *InMemoryLinkRepo) Save(_ context.Context, link *domain.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.slugs[link.Slug]; exists {
		return errors.New("slug already exists")
	}
	r.links[link.ID] = link
	r.slugs[link.Slug] = link.ID
	return nil
}

func (r *InMemoryLinkRepo) FindByID(_ context.Context, id domain.LinkID) (*domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	link, ok := r.links[id]
	if !ok {
		return nil, errs.ErrNotFound
	}
	return link, nil
}

func (r *InMemoryLinkRepo) FindBySlug(_ context.Context, slug string) (*domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.slugs[slug]
	if !ok {
		return nil, errs.ErrNotFound
	}
	return r.links[id], nil
}

func (r *InMemoryLinkRepo) FindAll(_ context.Context, filter domain.LinkFilter) ([]*domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	matches := make([]*domain.Link, 0, len(r.links))

	for _, link := range r.links {
		if filter.UserID != "" && link.UserID != filter.UserID {
			continue
		}

		if filter.CampaignID != "" && link.CampaignID != filter.CampaignID {
			continue
		}

		if filter.IsActive != nil && link.IsActive != *filter.IsActive {
			continue
		}

		if filter.Search != "" {
			term := strings.ToLower(filter.Search)
			slug := strings.ToLower(link.Slug)
			target := strings.ToLower(link.TargetURL)

			if !strings.Contains(slug, term) && !strings.Contains(target, term) {
				continue
			}
		}

		matches = append(matches, link)
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].CreatedAt.After(matches[j].CreatedAt)
	})

	if filter.Offset >= len(matches) {
		return []*domain.Link{}, nil
	}

	res := matches[filter.Offset:]

	if filter.Limit > 0 && filter.Limit < len(res) {
		res = res[:filter.Limit]
	}

	return res, nil
}

func (r *InMemoryLinkRepo) CountByUserID(_ context.Context, userID domain.UserID) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var count int64
	for _, l := range r.links {
		if l.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (r *InMemoryLinkRepo) Delete(_ context.Context, userID domain.UserID, id domain.LinkID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	link, ok := r.links[id]
	if !ok {
		return errs.ErrNotFound
	}

	if link.UserID != userID {
		return errs.ErrUnauthorized
	}

	delete(r.slugs, link.Slug)
	delete(r.links, id)

	return nil
}
