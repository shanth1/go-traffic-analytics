package cached

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type UserRepo struct {
	repo  ports.UserRepository
	cache ports.Cache
	ttl   time.Duration
}

func NewUserRepo(repo ports.UserRepository, cache ports.Cache, ttl time.Duration) ports.UserRepository {
	return &UserRepo{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

func (r *UserRepo) buildIDKey(id domain.UserID) string {
	return fmt.Sprintf("gotrace:user:id:%s", id)
}

func (r *UserRepo) buildEmailKey(email string) string {
	return fmt.Sprintf("gotrace:user:email:%s", email)
}

func (r *UserRepo) Save(ctx context.Context, user *domain.User) error {
	if err := r.repo.Save(ctx, user); err != nil {
		return err
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// TODO: logging
		_ = r.cache.Delete(bgCtx, r.buildIDKey(user.ID))
		_ = r.cache.Delete(bgCtx, r.buildEmailKey(user.Email))
	}()

	return nil
}

func (r *UserRepo) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	key := r.buildIDKey(id)

	val, err := r.cache.Get(ctx, key)
	if err == nil {
		if bytesVal, ok := val.([]byte); ok {
			var user domain.User
			if jsonErr := json.Unmarshal(bytesVal, &user); jsonErr == nil {
				return &user, nil
			}
			// TODO: logging (unmarshal error)
		}
	}
	// else if !errors.Is(err, ports.ErrCacheMiss) {
	// TODO: logging (cache error)
	// }

	user, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		bytes, err := json.Marshal(user)
		if err != nil {
			// TODO: logging (marshal error)
			return
		}

		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// TODO: logging (failed to set cache)
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return user, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	key := r.buildEmailKey(email)

	val, err := r.cache.Get(ctx, key)
	if err == nil {
		if bytesVal, ok := val.([]byte); ok {
			var user domain.User
			if jsonErr := json.Unmarshal(bytesVal, &user); jsonErr == nil {
				return &user, nil
			}
			// TODO: logging (unmarshal error)
		}
	}
	// else if !errors.Is(err, ports.ErrCacheMiss) {
	// TODO: logging (cache error)
	// }

	user, err := r.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	go func() {
		bytes, err := json.Marshal(user)
		if err != nil {
			// TODO: logging (marshal error)
			return
		}

		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// TODO: logging (failed to set cache)
		_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
	}()

	return user, nil
}

func (r *UserRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return r.repo.FindAll(ctx, limit, offset)
}

func (r *UserRepo) IncrementUsage(ctx context.Context, userID domain.UserID, delta int) error {
	if err := r.repo.IncrementUsage(ctx, userID, delta); err != nil {
		return err
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// TODO: logging
		_ = r.cache.Delete(bgCtx, r.buildIDKey(userID))
	}()

	return nil
}

// TODO:
// func (r *UserRepo) invalidate(ctx context.Context, id domain.UserID, email string) {
// 	go func() {
// 		bg := context.Background()
// 		_ = r.cache.Delete(bg, r.keyID(id))
// 		if email != "" {
// 			_ = r.cache.Delete(bg, r.keyEmail(email))
// 		}
// 	}()
// }
