package generator

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type DataSeeder struct {
	UserRepo     ports.UserRepository
	CampaignRepo ports.CampaignRepository
	LinkRepo     ports.LinkRepository
	ClickRepo    ports.ClickRepository
}

func (s *DataSeeder) SeedFullTopology() {
	ctx := context.Background()
	fmt.Println("🌱 Starting Full Topology Seeding...")

	// 1. Create Admin User
	admin := &domain.User{
		ID:           uuid.New().String(),
		Email:        "admin@tracebit.com",
		PasswordHash: "$2a$10$...", // mock hash for 'password'
		Role:         domain.RoleAdmin,
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
	_ = s.UserRepo.Save(ctx, admin)
	fmt.Printf("👤 Created Admin: %s\n", admin.Email)

	// 2. Create Client User
	client := &domain.User{
		ID:           uuid.New().String(),
		Email:        "client@example.com", // Этот email можно использовать для входа на фронте
		PasswordHash: "$2a$10$...",
		Role:         domain.RoleClient,
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
	_ = s.UserRepo.Save(ctx, client)
	fmt.Printf("👤 Created Client: %s\n", client.Email)

	// 3. Create Campaigns for Client
	campaignNames := []string{"Black Friday 2024", "Summer Sale", "Influencer Integrations"}

	for _, name := range campaignNames {
		campID := uuid.New().String()
		campaign := &domain.Campaign{
			ID:        campID,
			UserID:    client.ID,
			Name:      name,
			CreatedAt: time.Now(),
		}
		_ = s.CampaignRepo.Save(ctx, campaign)

		// 4. Create Links inside Campaign
		s.seedLinksForCampaign(ctx, campID)
	}

	fmt.Println("✅ Seeding completed!")
}

func (s *DataSeeder) seedLinksForCampaign(ctx context.Context, campaignID string) {
	// Генерация 3-5 ссылок на кампанию
	linksCount := rand.Intn(3) + 3

	for i := 0; i < linksCount; i++ {
		linkID := uuid.New().String()
		link := &domain.Link{
			ID:         linkID,
			CampaignID: campaignID,
			Slug:       fmt.Sprintf("lnk-%s-%d", campaignID[:4], i), // ex: lnk-a1b2-0
			TargetURL:  "https://google.com",
			IsActive:   true,
			CreatedAt:  time.Now(),
		}
		_ = s.LinkRepo.Save(ctx, link)

		// 5. Generate Traffic (Clicks)
		// Генерируем от 50 до 500 кликов на ссылку
		clicksCount := rand.Intn(450) + 50
		s.seedClicksForLink(ctx, linkID, clicksCount)
	}
}

func (s *DataSeeder) seedClicksForLink(ctx context.Context, linkID string, count int) {
	referers := []string{"Facebook", "Twitter", "Instagram", "Google Search", "Direct"}
	countries := []string{"US", "DE", "FR", "GB", "JP", "BR"}

	now := time.Now()

	for i := 0; i < count; i++ {
		daysAgo := rand.Intn(14)
		fakeTime := now.AddDate(0, 0, -daysAgo).Add(time.Duration(rand.Intn(24)) * time.Hour)

		_ = s.ClickRepo.Save(ctx, &domain.ClickEvent{
			ID:        uuid.New().String(),
			LinkID:    linkID,
			Timestamp: fakeTime,
			IP:        fmt.Sprintf("10.0.%d.%d", rand.Intn(255), rand.Intn(255)),
			Country:   countries[rand.Intn(len(countries))],
			OS:        "iOS",
			Browser:   "Safari",
			Device:    "Mobile",
			Referer:   referers[rand.Intn(len(referers))],
		})
	}
}
