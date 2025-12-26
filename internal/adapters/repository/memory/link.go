package memoryrepo

import (
	"context"
	"errors"
	"sync"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type InMemoryLinkRepo struct {
	mu    sync.RWMutex
	links map[string]*domain.Link
	slugs map[string]string // slug -> id mapping
}

func NewLinkRepo() ports.LinkRepository {
	return &InMemoryLinkRepo{
		links: make(map[string]*domain.Link),
		slugs: make(map[string]string),
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

func (r *InMemoryLinkRepo) FindBySlug(_ context.Context, slug string) (*domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.slugs[slug]
	if !ok {
		return nil, errors.New("not found")
	}
	return r.links[id], nil
}

func (r *InMemoryLinkRepo) FindAll(_ context.Context) ([]*domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Link, 0, len(r.links))
	for _, v := range r.links {
		result = append(result, v)
	}
	return result, nil
}

func (r *InMemoryLinkRepo) FindAllByCampaignID(_ context.Context, campaignID string) ([]*domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Link, 0, len(r.links))
	for _, v := range r.links {
		if v.CampaignID == campaignID {
			result = append(result, v)
		}
	}
	return result, nil
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
