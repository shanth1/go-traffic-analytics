package memoryrepo

import (
	"context"
	"fmt"
	"sync"

	"github.com/shanth1/gotools/errs"
	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type InMemoryUserRepo struct {
	mu     sync.RWMutex
	users  map[domain.UserID]*domain.User
	emails map[string]domain.UserID
}

func NewUserRepo() ports.UserRepository {
	return &InMemoryUserRepo{
		users:  make(map[domain.UserID]*domain.User),
		emails: make(map[string]domain.UserID),
	}
}

func (r *InMemoryUserRepo) Save(_ context.Context, user *domain.User) error {
	const op = "memoryrepo.InMemoryUserRepo.Save"

	r.mu.Lock()
	defer r.mu.Unlock()

	if existingID, exists := r.emails[user.Email]; exists {
		if existingID != user.ID {
			return ops.E(op, ops.KindExist, fmt.Errorf("email conflict: %q is owned by user %q (current user: %q): %w", user.Email, existingID, user.ID, errs.ErrAlreadyExists))
		}
	}

	r.users[user.ID] = user
	r.emails[user.Email] = user.ID

	return nil
}

func (r *InMemoryUserRepo) FindByID(_ context.Context, id domain.UserID) (*domain.User, error) {
	const op = "memoryrepo.InMemoryUserRepo.FindByID"

	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, ops.E(op, ops.KindNotFound, fmt.Errorf("user with id %q: %w", id, errs.ErrNotFound))
	}

	return u, nil
}

func (r *InMemoryUserRepo) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	const op = "memoryrepo.InMemoryUserRepo.FindByEmail"

	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.emails[email]
	if !ok {
		return nil, ops.E(op, ops.KindNotFound, fmt.Errorf("user with email %q: %w", email, errs.ErrNotFound))
	}

	return r.users[id], nil
}

func (r *InMemoryUserRepo) FindAll(_ context.Context, limit, offset int) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.User
	i := 0
	for _, u := range r.users {
		if i >= offset {
			if limit > 0 && len(result) >= limit {
				break
			}
			result = append(result, u)
		}
		i++
	}
	return result, nil
}

func (r *InMemoryUserRepo) Count(_ context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return int64(len(r.users)), nil
}

func (r *InMemoryUserRepo) IncrementUsage(_ context.Context, userID domain.UserID, delta int) error {
	const op = "memoryrepo.InMemoryUserRepo.IncrementUsage"

	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return ops.E(op, ops.KindNotFound, fmt.Errorf("user with id %q: %w", userID, errs.ErrNotFound))
	}

	user.ClicksCurrentMonth += delta
	return nil
}

func (r *InMemoryUserRepo) ResetUsage(_ context.Context, userID domain.UserID) error {
	const op = "memoryrepo.InMemoryUserRepo.ResetUsage"

	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return ops.E(op, ops.KindNotFound, fmt.Errorf("user with id %q: %w", userID, errs.ErrNotFound))
	}

	user.ClicksCurrentMonth = 0
	return nil
}
