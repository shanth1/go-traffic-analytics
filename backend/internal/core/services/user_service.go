package services

import (
	"context"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type UserService struct {
	userRepo ports.UserRepository
	planRepo ports.PlanRepository
}

func NewUserService(u ports.UserRepository, p ports.PlanRepository) *UserService {
	return &UserService{userRepo: u, planRepo: p}
}

func (s *UserService) GetAllUsers(ctx context.Context, page, limit int) ([]*domain.User, error) {
	offset := (page - 1) * limit
	return s.userRepo.FindAll(ctx, limit, offset)
}

func (s *UserService) SetUserStatus(ctx context.Context, userID string, isActive bool) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	user.IsActive = isActive
	return s.userRepo.Save(ctx, user)
}

func (s *UserService) ChangeUserPlan(ctx context.Context, userID, planID string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if _, err := s.planRepo.FindByID(ctx, planID); err != nil {
		return err
	}
	user.PlanID = planID
	return s.userRepo.Save(ctx, user)
}

func (s *UserService) GetAllPlans(ctx context.Context) ([]*domain.Plan, error) {
	return s.planRepo.FindAll(ctx)
}
