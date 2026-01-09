package cached

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shanth1/gotools/ops"
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
	const op = "cached.UserRepo.Save"

	if err := r.repo.Save(ctx, user); err != nil {
		return ops.E(op, err)
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
	const op = "cached.UserRepo.FindByID"

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
		return nil, ops.E(op, err)
	}

	go func() {
		bytes, err := json.Marshal(user)
		if err == nil {
			bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
		}
	}()

	return user, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	const op = "cached.UserRepo.FindByEmail"

	key := r.buildEmailKey(email)
	val, err := r.cache.Get(ctx, key)
	if err == nil {
		if bytesVal, ok := val.([]byte); ok {
			var user domain.User
			if jsonErr := json.Unmarshal(bytesVal, &user); jsonErr == nil {
				return &user, nil
			}
		}
	}

	user, err := r.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ops.E(op, err)
	}

	go func() {
		bytes, err := json.Marshal(user)
		if err == nil {
			bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = r.cache.Set(bgCtx, key, bytes, r.ttl)
		}
	}()

	return user, nil
}

func (r *UserRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	const op = "cached.UserRepo.FindAll"

	users, err := r.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, ops.E(op, err)
	}

	return users, nil
}

func (r *UserRepo) Count(ctx context.Context) (int64, error) {
	const op = "cached.UserRepo.Count"

	count, err := r.repo.Count(ctx)
	if err != nil {
		return 0, ops.E(op, err)
	}

	return count, nil
}

func (r *UserRepo) IncrementUsage(ctx context.Context, userID domain.UserID, delta int) error {
	const op = "cached.UserRepo.IncrementUsage"

	if err := r.repo.IncrementUsage(ctx, userID, delta); err != nil {
		return ops.E(op, err)
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.cache.Delete(bgCtx, r.buildIDKey(userID))
	}()

	return nil
}

func (r *UserRepo) ResetUsage(ctx context.Context, userID domain.UserID) error {
	const op = "cached.UserRepo.ResetUsage"

	if err := r.repo.ResetUsage(ctx, userID); err != nil {
		return ops.E(op, err)
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
