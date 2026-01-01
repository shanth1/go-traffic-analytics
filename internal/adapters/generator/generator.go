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

type Config struct {
	UsersCount       int
	CampaignPerUser  int
	LinksPerCampaign int
	ClicksPerLink    int
	DaysHistory      int
	PasswordDefault  string
}

type DataSeeder struct {
	userRepo      ports.UserRepository
	campaignRepo  ports.CampaignRepository
	linkRepo      ports.LinkRepository
	analyticsRepo ports.AnalyticsRepository
	planRepo      ports.PlanRepository
	rng           *rand.Rand
}

func New(
	ur ports.UserRepository,
	cr ports.CampaignRepository,
	lr ports.LinkRepository,
	ar ports.AnalyticsRepository,
	pr ports.PlanRepository,
) *DataSeeder {
	return &DataSeeder{
		userRepo:      ur,
		campaignRepo:  cr,
		linkRepo:      lr,
		analyticsRepo: ar,
		planRepo:      pr,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *DataSeeder) Seed(ctx context.Context, cfg Config) error {
	fmt.Println("🚀 Starting Advanced Data Seeding...")

	plans, err := s.seedPlans(ctx)
	if err != nil {
		return fmt.Errorf("seed plans: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.PasswordDefault), 10)
	if err != nil {
		return fmt.Errorf("password hash: %w", err)
	}
	admin := &domain.User{
		ID:           domain.UserID(uuid.NewString()),
		Email:        "admin@gotrace.com",
		PasswordHash: string(hash),
		Role:         domain.RoleAdmin,
		PlanID:       "enterprise",
		IsActive:     true,
		CreatedAt:    time.Now().AddDate(0, -1, 0),
	}
	if err := s.userRepo.Save(ctx, admin); err != nil {
		return err
	}

	for i := 0; i < cfg.UsersCount; i++ {
		plan := plans[s.rng.Intn(len(plans))]
		user := s.generateUser(fmt.Sprintf("%d", i), plan.ID, cfg.PasswordDefault)
		if err := s.userRepo.Save(ctx, user); err != nil {
			return err
		}

		for c := 0; c < s.rng.Intn(cfg.CampaignPerUser)+1; c++ {
			camp := s.generateCampaign(user.ID)
			if err := s.campaignRepo.Save(ctx, camp); err != nil {
				return err
			}

			for l := 0; l < cfg.LinksPerCampaign; l++ {
				link := s.generateLink(user.ID, camp.ID)
				if err := s.linkRepo.Save(ctx, link); err != nil {
					return err
				}

				if err := s.seedClicks(ctx, link.ID, link.CampaignID, cfg.ClicksPerLink, cfg.DaysHistory); err != nil {
					return err
				}
			}
		}
	}

	fmt.Println("✨ Seeding Completed Successfully!")
	return nil
}

func (s *DataSeeder) seedPlans(ctx context.Context) ([]domain.Plan, error) {
	plans := []domain.Plan{
		{
			ID:             "free",
			Name:           "Free Plan",
			Description:    "This is Free Plan",
			PriceCents:     0,
			MaxLinks:       5,
			MaxClicksMonth: 1000,
			CanExportData:  false,
			IsActive:       true,
		},
		{
			ID:             "pro",
			Name:           "Pro Plan",
			Description:    "This is Pro Plan",
			PriceCents:     1900,
			MaxLinks:       100,
			MaxClicksMonth: 50000,
			CanExportData:  true,
			IsActive:       true,
		},
		{
			ID:             "enterprise",
			Name:           "Enterprise",
			Description:    "This is Enterprise Plan",
			PriceCents:     9900,
			MaxLinks:       -1,
			MaxClicksMonth: 1000000,
			CanExportData:  true,
			IsActive:       true,
		},
	}
	for _, p := range plans {
		if err := s.planRepo.Save(ctx, &p); err != nil {
			return nil, err
		}
	}
	return plans, nil
}

func (s *DataSeeder) generateUser(id string, planID string, password string) *domain.User {
	email := fmt.Sprintf("user%s@example.com", id)
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 4)

	var expires *time.Time
	if s.rng.Float32() < 0.2 {
		t := time.Now().AddDate(0, 0, -5)
		expires = &t
	}

	return &domain.User{
		ID:           domain.UserID(id),
		Email:        email,
		PasswordHash: string(hash),
		Role:         domain.RoleClient,
		PlanID:       planID,
		PlanExpires:  expires,
		IsActive:     true,
		CreatedAt:    time.Now().AddDate(0, 0, -30),
	}
}

func (s *DataSeeder) generateCampaign(userID domain.UserID) *domain.Campaign {
	names := []string{"Google Ads", "TikTok Promo", "Email Newsletter", "Winter Sale", "Product Hunt Launch"}
	return &domain.Campaign{
		ID:        uuid.NewString(),
		UserID:    userID,
		Name:      names[s.rng.Intn(len(names))],
		CreatedAt: time.Now().AddDate(0, 0, -20),
	}
}

func (s *DataSeeder) generateLink(userID domain.UserID, campID string) *domain.Link {
	slugs := []string{"promo", "deal", "secret", "update", "buy-now", "discount"}
	return &domain.Link{
		ID:         domain.LinkID(uuid.NewString()),
		UserID:     userID,
		CampaignID: campID,
		Slug:       fmt.Sprintf("%s-%d", slugs[s.rng.Intn(len(slugs))], s.rng.Int63()),
		TargetURL:  "https://example.com/target",
		IsActive:   true,
		CreatedAt:  time.Now().AddDate(0, 0, -15),
	}
}

func (s *DataSeeder) seedClicks(ctx context.Context, linkID domain.LinkID, campID string, count int, days int) error {
	clickEvents := make([]*domain.ClickEvent, 0)
	for i := 0; i < count; i++ {
		hour := s.getWeightedHour()
		dayOffset := s.rng.Intn(days)
		timestamp := time.Now().AddDate(0, 0, -dayOffset).Truncate(time.Hour).Add(time.Duration(hour) * time.Hour)

		geo := s.getRandomGeo()
		ua := s.getRandomUserAgent()

		click := &domain.ClickEvent{
			ID:         uuid.NewString(),
			LinkID:     linkID,
			CampaignID: campID,
			Timestamp:  timestamp,
			IP:         fmt.Sprintf("%d.%d.%d.%d", s.rng.Intn(255), s.rng.Intn(255), s.rng.Intn(255), s.rng.Intn(255)),
			Country:    geo.Country,
			City:       geo.City,
			OS:         ua.OS,
			Browser:    ua.Browser,
			Device:     ua.Device,
			Referer:    s.getRandomReferer(),
		}
		clickEvents = append(clickEvents, click)
	}
	if err := s.analyticsRepo.SaveBatch(ctx, clickEvents); err != nil {
		return err
	}

	return nil
}

type geoPair struct{ Country, City string }

func (s *DataSeeder) getRandomGeo() geoPair {
	geos := []geoPair{
		{"US", "New York"}, {"US", "San Francisco"}, {"GB", "London"},
		{"DE", "Berlin"}, {"DE", "Munich"}, {"FR", "Paris"},
		{"RU", "Moscow"}, {"RU", "Saint Petersburg"}, {"JP", "Tokyo"},
	}
	return geos[s.rng.Intn(len(geos))]
}

type uaGroup struct{ OS, Browser, Device string }

func (s *DataSeeder) getRandomUserAgent() uaGroup {
	r := s.rng.Float32()
	switch {
	case r < 0.45:
		return uaGroup{"iOS", "Safari", "Mobile"}
	case r < 0.80:
		return uaGroup{"Windows", "Chrome", "Desktop"}
	case r < 0.90:
		return uaGroup{"Android", "Chrome", "Mobile"}
	case r < 0.95:
		return uaGroup{"macOS", "Safari", "Desktop"}
	default:
		return uaGroup{"Linux", "Firefox", "Desktop"}
	}
}

func (s *DataSeeder) getRandomReferer() string {
	refs := []string{"https://t.co", "https://facebook.com", "Direct", "https://google.com", "https://linkedin.com"}
	if s.rng.Float32() < 0.4 {
		return "Direct"
	}
	return refs[s.rng.Intn(len(refs))]
}

func (s *DataSeeder) getWeightedHour() int {
	h := s.rng.Intn(24)
	if h > 2 && h < 6 {
		if s.rng.Float32() < 0.8 {
			return s.rng.Intn(10) + 10
		}
	}
	return h
}
