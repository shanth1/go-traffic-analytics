package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/shanth1/gotools/errs"
	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type UserRepo struct {
	mu     sync.RWMutex
	users  map[domain.UserID]*domain.User
	emails map[string]domain.UserID
}

func NewUserRepo() ports.UserRepository {
	return &UserRepo{
		users:  make(map[domain.UserID]*domain.User),
		emails: make(map[string]domain.UserID),
	}
}

func (r *UserRepo) Save(_ context.Context, user *domain.User) error {
	const op = "memory.UserRepo.Save"

	r.mu.Lock()
	defer r.mu.Unlock()

	if existingID, exists := r.emails[user.Email]; exists {
		if existingID != user.ID {
			techErr := fmt.Errorf("email conflict: %q is owned by user %q (current user: %q): %w", user.Email, existingID, user.ID, errs.ErrAlreadyExists)
			return ops.WrapMsg(op, ops.KindExist, techErr, "email conflict")
		}
	}

	r.users[user.ID] = user
	r.emails[user.Email] = user.ID

	return nil
}

func (r *UserRepo) FindByID(_ context.Context, id domain.UserID) (*domain.User, error) {
	const op = "memory.UserRepo.FindByID"

	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, ops.WrapMsg(op, ops.KindNotFound, errs.ErrNotFound, fmt.Sprintf("%q user not found", id))
	}

	return u, nil
}

func (r *UserRepo) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	const op = "memory.UserRepo.FindByEmail"

	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.emails[email]
	if !ok {
		return nil, ops.WrapMsg(op, ops.KindNotFound, errs.ErrNotFound, fmt.Sprintf("user with email %q not found", email))
	}

	return r.users[id], nil
}

func (r *UserRepo) FindAll(_ context.Context, limit, offset int) ([]*domain.User, error) {
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

func (r *UserRepo) Count(_ context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return int64(len(r.users)), nil
}

func (r *UserRepo) IncrementUsage(_ context.Context, userID domain.UserID, delta int) error {
	const op = "memory.UserRepo.IncrementUsage"

	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return ops.WrapMsg(op, ops.KindNotFound, errs.ErrNotFound, fmt.Sprintf("%q user not found", userID))
	}

	user.ClicksCurrentMonth += delta
	return nil
}

func (r *UserRepo) ResetUsage(_ context.Context, userID domain.UserID) error {
	const op = "memory.UserRepo.ResetUsage"

	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return ops.WrapMsg(op, ops.KindNotFound, errs.ErrNotFound, fmt.Sprintf("%q user not found", userID))
	}

	user.ClicksCurrentMonth = 0
	return nil
}
