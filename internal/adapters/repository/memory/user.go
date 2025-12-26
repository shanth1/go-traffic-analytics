package memoryrepo

import (
	"context"
	"errors"
	"sync"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type InMemoryUserRepo struct {
	mu     sync.RWMutex
	users  map[domain.UserID]*domain.User
	emails map[string]domain.UserID // Email -> ID (index)
}

func NewUserRepo() ports.UserRepository {
	return &InMemoryUserRepo{
		users:  make(map[domain.UserID]*domain.User),
		emails: make(map[string]domain.UserID),
	}
}

func (r *InMemoryUserRepo) Save(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existingID, exists := r.emails[user.Email]; exists {
		if existingID != user.ID {
			return errors.New("email already exists")
		}
	}

	r.users[user.ID] = user
	r.emails[user.Email] = user.ID
	return nil
}

func (r *InMemoryUserRepo) FindByID(_ context.Context, id domain.UserID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (r *InMemoryUserRepo) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.emails[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return r.users[id], nil
}

func (r *InMemoryUserRepo) FindAll(_ context.Context, limit, offset int) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.User
	i := 0
	for _, u := range r.users {
		if i >= offset && len(result) < limit {
			result = append(result, u)
		}
		i++
	}
	return result, nil
}

func (r *InMemoryUserRepo) IncrementClickCount(_ context.Context, userID domain.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return errors.New("user not found")
	}

	user.ClicksCurrentMonth++
	return nil
}
