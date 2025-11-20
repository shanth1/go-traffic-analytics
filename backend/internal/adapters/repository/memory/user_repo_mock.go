package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type InMemoryUserRepo struct {
	mu     sync.RWMutex
	users  map[string]*domain.User // ID -> User
	emails map[string]string       // Email -> ID (индекс)
}

func NewUserRepo() ports.UserRepository {
	return &InMemoryUserRepo{
		users:  make(map[string]*domain.User),
		emails: make(map[string]string),
	}
}

func (r *InMemoryUserRepo) Save(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверка на уникальность email (если это создание нового юзера или смена email)
	if existingID, exists := r.emails[user.Email]; exists {
		if existingID != user.ID {
			return errors.New("email already exists")
		}
	}

	r.users[user.ID] = user
	r.emails[user.Email] = user.ID
	return nil
}

func (r *InMemoryUserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (r *InMemoryUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.emails[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return r.users[id], nil
}

func (r *InMemoryUserRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Мок пагинации (очень простой)
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
