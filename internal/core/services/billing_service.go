package services

import (
	"context"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type BillingService struct {
	planRepo ports.PlanRepository
}

func NewBillingService(p ports.PlanRepository) *BillingService {
	return &BillingService{planRepo: p}
}

func (s *BillingService) GetAllPlans(ctx context.Context) ([]*domain.Plan, error) {
	return s.planRepo.FindAll(ctx)
}
