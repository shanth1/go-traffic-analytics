package memoryrepo

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/shanth1/gotools/errs"
	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type MemoryLinkRepo struct {
	mu    sync.RWMutex
	links map[domain.LinkID]*domain.Link
	slugs map[string]domain.LinkID // slug -> linkID (index)
}

func NewLinkRepo() ports.LinkRepository {
	return &MemoryLinkRepo{
		links: make(map[domain.LinkID]*domain.Link),
		slugs: make(map[string]domain.LinkID),
	}
}

func (r *MemoryLinkRepo) Save(_ context.Context, link *domain.Link) error {
	const op = "memoryrepo.MemoryLinkRepo.Save"

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.slugs[link.Slug]; exists {
		return ops.E(op, ops.KindExist, fmt.Errorf("link with slug %q: %w", link.Slug, errs.ErrAlreadyExists))
	}
	r.links[link.ID] = link
	r.slugs[link.Slug] = link.ID
	return nil
}

func (r *MemoryLinkRepo) FindByID(_ context.Context, id domain.LinkID) (*domain.Link, error) {
	const op = "memoryrepo.MemoryLinkRepo.FindByID"

	r.mu.RLock()
	defer r.mu.RUnlock()

	link, ok := r.links[id]
	if !ok {
		return nil, ops.E(op, ops.KindNotFound, fmt.Errorf("link with id %q: %w", id, errs.ErrNotFound))
	}

	return link, nil
}

func (r *MemoryLinkRepo) FindBySlug(_ context.Context, slug string) (*domain.Link, error) {
	const op = "memoryrepo.MemoryLinkRepo.FindBySlug"

	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.slugs[slug]
	if !ok {
		return nil, ops.E(op, ops.KindNotFound, fmt.Errorf("link with slug %q: %w", slug, errs.ErrNotFound))
	}

	return r.links[id], nil
}

func (r *MemoryLinkRepo) matchesFilter(link *domain.Link, filter domain.LinkFilter) bool {
	if filter.UserID != "" && link.UserID != filter.UserID {
		return false
	}
	if filter.CampaignID != "" && link.CampaignID != filter.CampaignID {
		return false
	}
	if filter.IsActive != nil && link.IsActive != *filter.IsActive {
		return false
	}
	if filter.Search != "" {
		term := strings.ToLower(filter.Search)
		slug := strings.ToLower(link.Slug)
		target := strings.ToLower(link.TargetURL)
		if !strings.Contains(slug, term) && !strings.Contains(target, term) {
			return false
		}
	}
	return true
}

func (r *MemoryLinkRepo) FindAll(_ context.Context, filter domain.LinkFilter) ([]*domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	matches := make([]*domain.Link, 0)
	for _, link := range r.links {
		if r.matchesFilter(link, filter) {
			matches = append(matches, link)
		}
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

func (r *MemoryLinkRepo) Count(_ context.Context, filter domain.LinkFilter) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var count int64
	for _, link := range r.links {
		if r.matchesFilter(link, filter) {
			count++
		}
	}
	return count, nil
}

func (r *MemoryLinkRepo) Delete(_ context.Context, userID domain.UserID, id domain.LinkID) error {
	const op = "memoryrepo.MemoryLinkRepo.Delete"

	r.mu.Lock()
	defer r.mu.Unlock()

	link, ok := r.links[id]
	if !ok {
		return ops.E(op, ops.KindNotFound, fmt.Errorf("link with id %q: %w", id, errs.ErrNotFound))
	}

	if link.UserID != userID {
		return ops.E(op, ops.KindUnauthorized, fmt.Errorf("expected user with id %q, got %q: %w", userID, link.UserID, errs.ErrUnauthorized))
	}

	delete(r.slugs, link.Slug)
	delete(r.links, id)

	return nil
}
