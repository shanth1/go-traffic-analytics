package generator

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

type DataSeeder struct {
	UserRepo     ports.UserRepository
	CampaignRepo ports.CampaignRepository
	LinkRepo     ports.LinkRepository
	ClickRepo    ports.ClickRepository
	PlanRepo     ports.PlanRepository
}

func (s *DataSeeder) SeedFullTopology() {
	ctx := context.Background()
	fmt.Println("🌱 Starting Full Topology Seeding...")

	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	realPasswordHash := string(hash)

	admin := &domain.User{
		ID:           uuid.New().String(),
		Email:        "admin@gotrace.com",
		PasswordHash: realPasswordHash,
		Role:         domain.RoleAdmin,
		IsActive:     true,
		PlanID:       "enterprise",
		CreatedAt:    time.Now(),
	}
	_ = s.UserRepo.Save(ctx, admin)
	fmt.Printf("👤 Created Admin: %s (password: password)\n", admin.Email)

	client := &domain.User{
		ID:           uuid.New().String(),
		Email:        "client@gotrace.com",
		PasswordHash: realPasswordHash,
		Role:         domain.RoleClient,
		PlanID:       "pro",
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
	_ = s.UserRepo.Save(ctx, client)
	fmt.Printf("👤 Created Client: %s (password: password)\n", client.Email)

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

		s.seedLinksForCampaign(ctx, client.ID, campID)
	}

	testLink := &domain.Link{
		ID: uuid.NewString(), UserID: client.ID, CampaignID: "",
		Slug: "google", TargetURL: "https://google.com", IsActive: true, CreatedAt: time.Now(),
	}
	_ = s.LinkRepo.Save(ctx, testLink)
	fmt.Println("🔗 Created Manual Link: /google -> https://google.com")

	fmt.Println("✅ Seeding completed!")
}

func (s *DataSeeder) seedLinksForCampaign(ctx context.Context, userID, campaignID string) {
	linksCount := rand.Intn(3) + 3

	for i := 0; i < linksCount; i++ {
		linkID := uuid.New().String()
		link := &domain.Link{
			ID:         linkID,
			UserID:     userID,
			CampaignID: campaignID,
			Slug:       fmt.Sprintf("lnk-%s-%d", campaignID[:4], i),
			TargetURL:  "https://google.com",
			IsActive:   true,
			CreatedAt:  time.Now(),
		}
		_ = s.LinkRepo.Save(ctx, link)

		clicksCount := rand.Intn(50) + 10
		s.seedClicksForLink(ctx, linkID, clicksCount)
	}
}

func (s *DataSeeder) seedClicksForLink(ctx context.Context, linkID string, count int) {
	referers := []string{"Facebook", "Twitter", "Instagram", "Google Search", "Direct"}
	countries := []string{"US", "DE", "FR", "GB", "JP", "BR"}

	now := time.Now()

	for i := 0; i < count; i++ {
		daysAgo := rand.Intn(7)
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
