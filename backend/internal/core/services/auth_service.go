package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type AuthService struct {
	userRepo ports.UserRepository
	planRepo ports.PlanRepository
}

func NewAuthService(u ports.UserRepository, p ports.PlanRepository) *AuthService {
	return &AuthService{userRepo: u, planRepo: p}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	// TODO:
	// 1. Проверка email (есть ли уже такой)

	// 2. Получаем дефолтный тариф
	defaultPlan, err := s.planRepo.FindDefault(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Создаем юзера
	user := &domain.User{
		ID:        uuid.New().String(),
		Email:     email,
		Role:      domain.RoleClient,
		PlanID:    defaultPlan.ID,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	// 4. Хэширование пароля и сохранение...
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
