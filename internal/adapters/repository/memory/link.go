package memory

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/shanth1/gotools/errs"
	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type LinkRepo struct {
	mu    sync.RWMutex
	links map[domain.LinkID]*domain.Link
	slugs map[string]domain.LinkID // slug -> linkID (index)
}

func NewLinkRepo() ports.LinkRepository {
	return &LinkRepo{
		links: make(map[domain.LinkID]*domain.Link),
		slugs: make(map[string]domain.LinkID),
	}
}

func (r *LinkRepo) Save(_ context.Context, link *domain.Link) error {
	const op = "memory.LinkRepo.Save"

	r.mu.Lock()
	defer r.mu.Unlock()

	if existingID, exists := r.slugs[link.Slug]; exists {
		if existingID != link.ID {
			return ops.WrapMsg(op, ops.KindExist, errs.ErrAlreadyExists, fmt.Sprintf("link with slug %q already exists", link.Slug))
		}
	}

	r.links[link.ID] = link
	r.slugs[link.Slug] = link.ID
	return nil
}

func (r *LinkRepo) FindByID(_ context.Context, id domain.LinkID) (*domain.Link, error) {
	const op = "memory.LinkRepo.FindByID"

	r.mu.RLock()
	defer r.mu.RUnlock()

	link, ok := r.links[id]
	if !ok || link.DeletedAt != nil {
		return nil, ops.WrapMsg(op, ops.KindNotFound, errs.ErrNotFound, fmt.Sprintf("%q link not found", id))
	}

	return link, nil
}

func (r *LinkRepo) FindBySlug(_ context.Context, slug string) (*domain.Link, error) {
	const op = "memory.LinkRepo.FindBySlug"

	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.slugs[slug]
	if !ok {
		return nil, ops.WrapMsg(op, ops.KindNotFound, errs.ErrNotFound, fmt.Sprintf("link with slug %q not found", slug))
	}

	link := r.links[id]
	if link.DeletedAt != nil {
		return nil, ops.WrapMsg(op, ops.KindNotFound, errs.ErrNotFound, fmt.Sprintf("link with slug %q not found", slug))
	}

	return link, nil
}

func (r *LinkRepo) matchesFilter(link *domain.Link, filter domain.LinkFilter) bool {
	if link.DeletedAt != nil {
		return false
	}
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

func (r *LinkRepo) FindAll(_ context.Context, filter domain.LinkFilter) ([]*domain.Link, error) {
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

func (r *LinkRepo) Count(_ context.Context, filter domain.LinkFilter) (int64, error) {
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

func (r *LinkRepo) Delete(_ context.Context, userID domain.UserID, id domain.LinkID) error {
	const op = "memory.LinkRepo.Delete"

	r.mu.Lock()
	defer r.mu.Unlock()

	link, ok := r.links[id]
	if !ok || link.DeletedAt != nil {
		return ops.WrapMsg(op, ops.KindNotFound, errs.ErrNotFound, fmt.Sprintf("%q link not found", id))
	}

	if link.UserID != userID {
		techErr := fmt.Errorf("user %q attempted to access link %q owned by %q: %w",
			userID, id, link.UserID, errs.ErrForbidden)
		return ops.WrapMsg(op, ops.KindPermission, techErr, "access denied")
	}

	now := time.Now()
	link.DeletedAt = &now

	return nil
}

func (r *LinkRepo) DeleteByUserID(_ context.Context, userID domain.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for _, link := range r.links {
		if link.UserID == userID && link.DeletedAt == nil {
			link.DeletedAt = &now
		}
	}
	return nil
}

func (r *LinkRepo) DeleteByCampaignID(_ context.Context, campaignID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for _, link := range r.links {
		if link.CampaignID == campaignID && link.DeletedAt == nil {
			link.DeletedAt = &now
		}
	}
	return nil
}
