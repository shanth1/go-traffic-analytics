package cachedproxy

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type UserRepo struct {
	repo   ports.UserRepository
	cache  ports.Cache
	logger log.Logger
	ttl    time.Duration
}

func NewUserRepo(repo ports.UserRepository, cache ports.Cache, logger log.Logger, ttl time.Duration) ports.UserRepository {
	return &UserRepo{
		repo:   repo,
		cache:  cache,
		logger: logger,
		ttl:    ttl,
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
	r.invalidate(user.ID, user.Email)
	return nil
}

func (r *UserRepo) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	key := r.buildIDKey(id)

	if user, err := r.getFromCache(ctx, key); err == nil && user != nil {
		return user, nil
	}

	user, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	r.setCache(key, user)
	return user, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	key := r.buildEmailKey(email)

	if user, err := r.getFromCache(ctx, key); err == nil && user != nil {
		return user, nil
	}

	user, err := r.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	r.setCache(key, user)
	return user, nil
}

func (r *UserRepo) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return r.repo.FindAll(ctx, limit, offset)
}

func (r *UserRepo) Count(ctx context.Context) (int64, error) {
	return r.repo.Count(ctx)
}

func (r *UserRepo) IncrementUsage(ctx context.Context, userID domain.UserID, delta int) error {
	if err := r.repo.IncrementUsage(ctx, userID, delta); err != nil {
		return err
	}
	r.invalidate(userID, "")
	return nil
}

func (r *UserRepo) ResetUsage(ctx context.Context, userID domain.UserID) error {
	if err := r.repo.ResetUsage(ctx, userID); err != nil {
		return err
	}
	r.invalidate(userID, "")
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id domain.UserID) error {
	var emailToDelete string

	keyID := r.buildIDKey(id)
	if cachedUser, _ := r.getFromCache(ctx, keyID); cachedUser != nil {
		emailToDelete = cachedUser.Email
	} else {
		if u, err := r.repo.FindByID(ctx, id); err == nil {
			emailToDelete = u.Email
		}
	}

	if err := r.repo.Delete(ctx, id); err != nil {
		return err
	}

	r.invalidate(id, emailToDelete)
	return nil
}

func (r *UserRepo) getFromCache(ctx context.Context, key string) (*domain.User, error) {
	val, err := r.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	bytesVal, ok := val.([]byte)
	if !ok {
		return nil, nil
	}
	var user domain.User
	if err := json.Unmarshal(bytesVal, &user); err != nil {
		r.logger.Error().Err(err).Str("key", key).Msg("unmarshal error in user repo")
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) setCache(key string, user *domain.User) {
	go func() {
		bytes, err := json.Marshal(user)
		if err != nil {
			r.logger.Error().Err(err).Msg("marshal error")
			return
		}
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := r.cache.Set(bgCtx, key, bytes, r.ttl); err != nil {
			r.logger.Error().Err(err).Str("key", key).Msg("failed to set cache")
		}
	}()
}

func (r *UserRepo) invalidate(id domain.UserID, email string) {
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.cache.Delete(bgCtx, r.buildIDKey(id)); err != nil {
			r.logger.Error().Err(err).Str("userID", string(id)).Msg("failed to invalidate user id")
		}

		if email != "" {
			if err := r.cache.Delete(bgCtx, r.buildEmailKey(email)); err != nil {
				r.logger.Error().Err(err).Str("email", email).Msg("failed to invalidate user email")
			}
		}
	}()
}
