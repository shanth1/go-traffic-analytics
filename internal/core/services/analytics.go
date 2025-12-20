package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/consts"
)

type AnalyticsService struct {
	clickRepo ports.ClickRepository
}

func NewAnalyticsService(c ports.ClickRepository) *AnalyticsService {
	return &AnalyticsService{clickRepo: c}
}

func (s *AnalyticsService) GetSummary(ctx context.Context, filter ports.AnalyticsFilter) (*domain.Summary, error) {
	total, err := s.clickRepo.CountTotal(ctx, filter)
	if err != nil {
		return nil, err
	}

	browsers, err := s.clickRepo.GetTopStats(ctx, filter, consts.Browser, 5)
	if err != nil {
		// TODO: logging
		browsers = []domain.CategoryStat{}
	}

	osStats, err := s.clickRepo.GetTopStats(ctx, filter, consts.OS, 5)
	if err != nil {
		// TODO: logging
		osStats = []domain.CategoryStat{}
	}

	// TODO: Unique Users, Top Country etc.
	return &domain.Summary{
		TotalClicks: total,
		TopBrowsers: browsers,
		TopOS:       osStats,
	}, nil
}

func (s *AnalyticsService) GetStreamGraphData(ctx context.Context, filter ports.AnalyticsFilter, groupBy string) ([]domain.StackedPoint, error) {
	duration := filter.To.Sub(filter.From)
	interval := time.Hour * 24 // Day by default

	if duration < time.Hour*48 {
		interval = time.Hour
	}

	return s.clickRepo.GetTimeSeriesGrouped(ctx, filter, groupBy, interval)
}

func (s *AnalyticsService) GetSankeyData(ctx context.Context, filter ports.AnalyticsFilter, stages []string) (*domain.SankeyData, error) {
	if len(stages) == 0 {
		stages = []string{consts.Referer, consts.Device, consts.Country}
	}

	return s.clickRepo.GetFlowData(ctx, filter, stages)
}

func (s *AnalyticsService) GetGeoDistribution(ctx context.Context, filter ports.AnalyticsFilter) (interface{}, error) {
	stats, err := s.clickRepo.GetTopStats(ctx, filter, consts.Country, 200)
	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, stat := range stats {
		result[stat.Name] = stat.Value
	}
	return result, nil
}

func (s *AnalyticsService) GetHeatmapData(ctx context.Context, filter ports.AnalyticsFilter) ([]domain.HeatmapPoint, error) {
	return s.clickRepo.GetHeatmapData(ctx, filter)
}

func (s *AnalyticsService) GetCategoryStats(ctx context.Context, filter ports.AnalyticsFilter, dimension string) ([]domain.CategoryStat, error) {
	return s.clickRepo.GetTopStats(ctx, filter, dimension, 10)
}

func (s *AnalyticsService) GetTrafficQuality(ctx context.Context, filter ports.AnalyticsFilter) (*domain.TrafficQuality, error) {
	total, err := s.clickRepo.CountTotal(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("counting total clicks: %w", err)
	}

	if total == 0 {
		return &domain.TrafficQuality{
			HumanScore: 100,
		}, nil
	}

	// 1. Mobile Friendly
	devices, err := s.clickRepo.GetTopStats(ctx, filter, consts.Device, 100)
	if err != nil {
		return nil, fmt.Errorf("get top device stats: %w", err)
	}
	mobileClicks := 0
	for _, d := range devices {
		if d.Name == "Mobile" || d.Name == "Tablet" {
			mobileClicks += d.Value
		}
	}
	mobileScore := int((float64(mobileClicks) / float64(total)) * 100)

	// 2. Bot score
	// А. OS check
	osStats, err := s.clickRepo.GetTopStats(ctx, filter, consts.OS, 100)
	if err != nil {
		return nil, fmt.Errorf("get top os stats: %w", err)
	}
	suspiciousOsClicks := 0
	badOS := map[string]bool{
		consts.Unknown: true,
		"Bot":          true,
		"Linux":        true,
		"Unix":         true,
	}

	for _, stat := range osStats {
		if badOS[stat.Name] {
			suspiciousOsClicks += stat.Value
		}
	}

	// B. Browser check
	browserStats, err := s.clickRepo.GetTopStats(ctx, filter, consts.Browser, 100)
	if err != nil {
		return nil, fmt.Errorf("get top browser stats: %w", err)
	}
	suspiciousBrowserClicks := 0
	badBrowsers := map[string]bool{
		consts.Unknown:      true,
		"Bot":               true,
		"curl":              true,
		"python-requests":   true,
		"Go-http-client":    true,
		"HeadlessChrome":    true,
		"Apache-HttpClient": true,
	}

	for _, stat := range browserStats {
		if stat.Name == "" || badBrowsers[stat.Name] {
			suspiciousBrowserClicks += stat.Value
		}
	}

	pctBadOS := (float64(suspiciousOsClicks) / float64(total)) * 100
	pctBadBrowser := (float64(suspiciousBrowserClicks) / float64(total)) * 100

	botScore := int(math.Max(pctBadOS, pctBadBrowser))
	if botScore > 100 {
		botScore = 100
	}

	// 3. Geo Diversity
	countries, err := s.clickRepo.GetTopStats(ctx, filter, consts.Country, 100)
	if err != nil {
		return nil, fmt.Errorf("get top county stats: %w", err)
	}
	geoScore := 0
	if len(countries) > 0 {
		topCountryClicks := countries[0].Value
		concentration := float64(topCountryClicks) / float64(total)

		geoScore = int((1.0 - concentration) * 100)
	}

	return &domain.TrafficQuality{
		MobileFriendlyScore: mobileScore,
		BotScore:            botScore,
		HumanScore:          100 - botScore,
		GeoDiversityScore:   geoScore,
		IsSuspicious:        false,
	}, nil
}
