package cached

import (
	"context"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type CachedUserRepo struct {
	repo  ports.UserRepository
	cache ports.Cache
	ttl   time.Duration
}

func NewUserRepo(repo ports.UserRepository, cache ports.Cache, ttl time.Duration) ports.UserRepository {
	return &CachedUserRepo{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (r *CachedUserRepo) Save(_ context.Context, user *domain.User) error {

}

func (r *CachedUserRepo) FindByID(_ context.Context, id string) (*domain.User, error) {

}

func (r *CachedUserRepo) FindByEmail(_ context.Context, email string) (*domain.User, error) {

}

func (r *CachedUserRepo) FindAll(_ context.Context, limit, offset int) ([]*domain.User, error) {

}

func (r *CachedUserRepo) IncrementClickCount(_ context.Context, userID string) error {

}
