package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shanth1/gotrace/internal/config"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  ports.UserRepository
	planRepo  ports.PlanRepository
	jwtSecret []byte
}

func NewAuthService(u ports.UserRepository, p ports.PlanRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:  u,
		planRepo:  p,
		jwtSecret: []byte(cfg.JWTSecret),
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	if _, err := s.userRepo.FindByEmail(ctx, email); err == nil {
		return nil, errors.New("email already registered")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	defaultPlan, err := s.planRepo.FindDefault(ctx)
	if err != nil {
		return nil, errors.New("default plan not configured")
	}

	user := &domain.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashedBytes),
		Role:         domain.RoleClient,
		PlanID:       defaultPlan.ID,
		IsActive:     true,
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", nil, errors.New("account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	claims := domain.JwtCustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return t, user, nil
}
