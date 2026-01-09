package generator

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

type DataSeeder struct {
	logger        log.Logger
	userRepo      ports.UserRepository
	campaignRepo  ports.CampaignRepository
	linkRepo      ports.LinkRepository
	analyticsRepo ports.AnalyticsRepository
	planRepo      ports.PlanRepository
	rng           *rand.Rand
}

func New(
	logger log.Logger,
	ur ports.UserRepository,
	cr ports.CampaignRepository,
	lr ports.LinkRepository,
	ar ports.AnalyticsRepository,
	pr ports.PlanRepository,
) *DataSeeder {
	logger = logger.With(log.Str("module", "generator"))
	return &DataSeeder{
		logger:        logger,
		userRepo:      ur,
		campaignRepo:  cr,
		linkRepo:      lr,
		analyticsRepo: ar,
		planRepo:      pr,
	}
}

func (s *DataSeeder) Seed(ctx context.Context, cfg Config) error {
	var seed = cfg.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	s.rng = rand.New(rand.NewSource(seed))

	s.logger.Info().Int64("seed", seed).Msg("starting data seeding")

	plans, err := s.seedPlans(ctx)
	if err != nil {
		return fmt.Errorf("seed plans: %w", err)
	}

	defaultPassHash, err := bcrypt.GenerateFromPassword([]byte(cfg.PasswordDefault), 10)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	passHashStr := string(defaultPassHash)

	if err := s.createAdmin(ctx, passHashStr); err != nil {
		return fmt.Errorf("create admin: %w", err)
	}

	for i := 1; i <= cfg.UsersCount; i++ {
		var (
			userTotalLinks  = 0
			userTotalClicks = 0
		)

		plan := plans[s.rng.Intn(len(plans))]
		user := s.generateUser(i, plan.ID, passHashStr)

		if err := s.userRepo.Save(ctx, user); err != nil {
			return fmt.Errorf("save user %d: %w", i, err)
		}

		campCount := s.randomRange(cfg.CampaignsPerUser.Min, cfg.CampaignsPerUser.Max)

		for c := 0; c < campCount; c++ {
			camp := s.generateCampaign(user.ID)
			if err := s.campaignRepo.Save(ctx, camp); err != nil {
				return fmt.Errorf("save campaign: %w", err)
			}

			linkCount := s.randomRange(cfg.LinksPerCampaign.Min, cfg.LinksPerCampaign.Max)
			userTotalLinks += linkCount

			for l := 0; l < linkCount; l++ {
				link := s.generateLink(user.ID, camp.ID)
				if err := s.linkRepo.Save(ctx, link); err != nil {
					return fmt.Errorf("save link: %w", err)
				}

				clicksCount := s.randomRange(cfg.ClicksPerLink.Min, cfg.ClicksPerLink.Max)
				isViral := s.rng.Float64() < cfg.ViralLinkProbability
				if isViral {
					clicksCount *= 10
					s.logger.Debug().Str("slug", link.Slug).Int("clicks", clicksCount).Msg("Viral link generated")
				}

				userTotalClicks += clicksCount

				if err := s.seedClicksBatch(ctx, link, clicksCount, cfg); err != nil {
					return fmt.Errorf("seed clicks batch: %w", err)
				}
			}
		}

		s.logger.Info().
			Str("plan", plan.Name).
			Str("email", user.Email).
			Int("campaigns", campCount).
			Int("links", userTotalLinks).
			Int("clicks", userTotalClicks).
			Msg("user seeded")
	}

	s.logger.Info().Msg("seeding completed successfully")
	return nil
}

func (s *DataSeeder) seedClicksBatch(ctx context.Context, link *domain.Link, totalClicks int, cfg Config) error {
	batch := make([]*domain.ClickEvent, 0, cfg.BatchSize)

	targetGeo := s.pickWeightedGeo()

	now := time.Now()
	startTime := now.AddDate(0, 0, -cfg.HistoryDays)

	for i := 0; i < totalClicks; i++ {
		ts := s.generateSmartTimestamp(startTime, now)

		geo := targetGeo
		if s.rng.Float64() > 0.7 {
			geo = s.pickWeightedGeo()
		}

		city := geo.Cities[s.rng.Intn(len(geo.Cities))]
		device := deviceProfiles[s.rng.Intn(len(deviceProfiles))]
		referer := referers[s.rng.Intn(len(referers))]

		click := &domain.ClickEvent{
			ID:         uuid.NewString(),
			LinkID:     link.ID,
			CampaignID: link.CampaignID,
			UserID:     link.UserID,
			Timestamp:  ts,
			IP:         s.generateRealisticIP(),
			Country:    geo.CountryCode,
			City:       city,
			OS:         device.OS,
			Browser:    device.Browser,
			Device:     device.Device,
			Referer:    referer,
		}

		batch = append(batch, click)

		if len(batch) >= cfg.BatchSize {
			if err := s.analyticsRepo.SaveBatch(ctx, batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if err := s.analyticsRepo.SaveBatch(ctx, batch); err != nil {
			return err
		}
	}

	return nil
}

// --- Generators Implementation ---

func (s *DataSeeder) seedPlans(ctx context.Context) ([]domain.Plan, error) {
	plans := []domain.Plan{
		{ID: "free", Name: "Free", PriceCents: 0, MaxLinks: 5, MaxClicksMonth: 1000, CanExportData: false, IsActive: true},
		{ID: "pro", Name: "Pro", PriceCents: 1900, MaxLinks: 100, MaxClicksMonth: 50000, CanExportData: true, IsActive: true},
		{ID: "enterprise", Name: "Enterprise", PriceCents: 9900, MaxLinks: -1, MaxClicksMonth: 1000000, CanExportData: true, IsActive: true},
	}
	for _, p := range plans {
		// TODO: return error without duplicate key error
		_ = s.planRepo.Save(ctx, &p)
	}
	return plans, nil
}

func (s *DataSeeder) createAdmin(ctx context.Context, hash string) error {
	email := "admin@gotrace.com"
	admin := &domain.User{
		ID:           domain.UserID(uuid.NewString()),
		Email:        email,
		PasswordHash: hash,
		Role:         domain.RoleAdmin,
		PlanID:       "enterprise",
		IsActive:     true,
		CreatedAt:    time.Now().AddDate(0, -6, 0),
	}

	if err := s.userRepo.Save(ctx, admin); err != nil {
		return err
	}

	s.logger.Info().Str("email", email).Msg("admin created")
	return nil
}

func (s *DataSeeder) generateUser(index int, planID string, hash string) *domain.User {
	return &domain.User{
		ID:           domain.UserID(uuid.NewString()),
		Email:        fmt.Sprintf("user%d@example.com", index),
		PasswordHash: hash,
		Role:         domain.RoleClient,
		PlanID:       planID,
		IsActive:     true,
		CreatedAt:    time.Now().AddDate(0, 0, -s.rng.Intn(60)),
	}
}

func (s *DataSeeder) generateCampaign(userID domain.UserID) *domain.Campaign {
	name := campaignNames[s.rng.Intn(len(campaignNames))]
	if s.rng.Float32() > 0.5 {
		name = fmt.Sprintf("%s %d", name, time.Now().Year()-s.rng.Intn(3))
	}
	return &domain.Campaign{
		ID:        uuid.NewString(),
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now().AddDate(0, 0, -30),
	}
}

func (s *DataSeeder) generateLink(userID domain.UserID, campID string) *domain.Link {
	slugWord := slugWords[s.rng.Intn(len(slugWords))]
	return &domain.Link{
		ID:         domain.LinkID(uuid.NewString()),
		UserID:     userID,
		CampaignID: campID,
		Slug:       fmt.Sprintf("%s-%s", slugWord, s.randomString(4)),
		TargetURL:  fmt.Sprintf("https://example.com/products/%d", s.rng.Intn(1000)),
		IsActive:   true,
		CreatedAt:  time.Now().AddDate(0, 0, -20),
	}
}

// --- Helpers & Distributions ---

func (s *DataSeeder) generateSmartTimestamp(start, end time.Time) time.Time {
	span := end.Sub(start)

	randomDuration := time.Duration(s.rng.Int63n(int64(span)))
	baseTime := start.Add(randomDuration)

	hour := (s.rng.Intn(24) + s.rng.Intn(24) + s.rng.Intn(24)) / 3

	y, m, d := baseTime.Date()
	finalTime := time.Date(y, m, d, hour, s.rng.Intn(60), s.rng.Intn(60), 0, baseTime.Location())

	if finalTime.After(end) {
		return end.Add(-time.Minute)
	}
	return finalTime
}

func (s *DataSeeder) pickWeightedGeo() GeoLocation {
	totalWeight := 0
	for _, w := range worldData {
		totalWeight += w.Weight
	}

	r := s.rng.Intn(totalWeight)
	current := 0
	for _, w := range worldData {
		current += w.Weight
		if r < current {
			return w.Geo
		}
	}
	return worldData[0].Geo // Fallback
}

func (s *DataSeeder) generateRealisticIP() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		s.rng.Intn(220)+1,
		s.rng.Intn(255),
		s.rng.Intn(255),
		s.rng.Intn(254)+1)
}

func (s *DataSeeder) randomRange(minVal, maxVal int) int {
	if minVal >= maxVal {
		return minVal
	}
	return s.rng.Intn(maxVal-minVal+1) + minVal
}

func (s *DataSeeder) randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[s.rng.Intn(len(letters))]
	}
	return string(b)
}
