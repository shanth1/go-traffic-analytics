// internal/core/services/link_service.go
package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

var ErrLimitReached = errors.New("plan limit reached: upgrade your plan to create more links")

type LinkService struct {
	linkRepo ports.LinkRepository
	userRepo ports.UserRepository
	planRepo ports.PlanRepository
}

func NewLinkService(l ports.LinkRepository, u ports.UserRepository, p ports.PlanRepository) *LinkService {
	return &LinkService{linkRepo: l, userRepo: u, planRepo: p}
}

func (s *LinkService) CreateLink(ctx context.Context, userID string, targetURL string, customSlug string) (*domain.Link, error) {
	// 1. Получаем пользователя, чтобы узнать PlanID
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. Получаем детали тарифа
	plan, err := s.planRepo.FindByID(ctx, user.PlanID)
	if err != nil {
		return nil, err
	}

	// 3. ПРОВЕРКА ЛИМИТОВ (Business Logic)
	if plan.MaxLinks != -1 { // -1 означает безлимит
		currentCount, err := s.linkRepo.CountByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}

		if int(currentCount) >= plan.MaxLinks {
			return nil, ErrLimitReached
		}
	}

	// 4. Создание ссылки (если лимиты ок)
	link := &domain.Link{
		ID: uuid.New().String(),
		// CampaignID: // TODO
		TargetURL: targetURL,
		Slug:      customSlug, // Или генерация рандомного
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := s.linkRepo.Save(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}
