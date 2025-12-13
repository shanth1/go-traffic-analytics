package services

import (
	"context"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type AnalyticsService struct {
	clickRepo ports.ClickRepository
}

func NewAnalyticsService(c ports.ClickRepository) *AnalyticsService {
	return &AnalyticsService{clickRepo: c}
}

func (s *AnalyticsService) GetSummary(ctx context.Context, filter ports.AnalyticsFilter) (map[string]interface{}, error) {
	total, err := s.clickRepo.CountTotal(ctx, filter)
	if err != nil {
		return nil, err
	}

	browsers, err := s.clickRepo.GetTopStats(ctx, filter, "browser", 5)
	if err != nil {
		// TODO: logging
		browsers = []domain.CategoryStat{}
	}

	osStats, err := s.clickRepo.GetTopStats(ctx, filter, "os", 5)
	if err != nil {
		// TODO: logging
		osStats = []domain.CategoryStat{}
	}

	// TODO: Unique Users, Top Country etc.
	return map[string]interface{}{
		"total_clicks": total,
		"top_browsers": browsers,
		"top_os":       osStats,
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
		stages = []string{"referer", "device", "country"}
	}

	return s.clickRepo.GetFlowData(ctx, filter, stages)
}

func (s *AnalyticsService) GetGeoDistribution(ctx context.Context, filter ports.AnalyticsFilter) (interface{}, error) {
	stats, err := s.clickRepo.GetTopStats(ctx, filter, "country", 200)
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

// GetTrafficQuality - Рассчитывает метрики качества трафика
func (s *AnalyticsService) GetTrafficQuality(ctx context.Context, filter ports.AnalyticsFilter) (map[string]int, error) {
	total, err := s.clickRepo.CountTotal(ctx, filter)
	if err != nil || total == 0 {
		// Возвращаем дефолтные значения, если данных нет
		return map[string]int{
			"mobile_friendly": 0,
			"bot_score":       0,
			"unique_ip":       0, // Сложно посчитать без raw queries
			"human_score":     100,
		}, nil
	}

	// 1. Считаем Mobile Friendly (% мобильных устройств)
	devices, _ := s.clickRepo.GetTopStats(ctx, filter, "device", 100)
	mobileClicks := 0
	for _, d := range devices {
		if d.Name == "Mobile" || d.Name == "Tablet" {
			mobileClicks += d.Value
		}
	}
	mobileScore := int((float64(mobileClicks) / float64(total)) * 100)

	// 2. Считаем Bot Score (на основе "Unknown" OS или Browser)
	// Это эвристика. В реальности нужны спец. базы ботов.
	osStats, _ := s.clickRepo.GetTopStats(ctx, filter, "os", 100)
	suspiciousClicks := 0
	for _, stat := range osStats {
		if stat.Name == "Unknown" || stat.Name == "Bot" {
			suspiciousClicks += stat.Value
		}
	}
	botScore := int((float64(suspiciousClicks) / float64(total)) * 100)

	// 3. Geo Diversity (простая метрика: если стран > 5, то 100%, иначе пропорционально)
	countries, _ := s.clickRepo.GetTopStats(ctx, filter, "country", 100)
	geoScore := len(countries) * 20
	if geoScore > 100 {
		geoScore = 100
	}

	return map[string]int{
		"mobile_friendly": mobileScore,
		"bot_score":       botScore,
		"human_score":     100 - botScore,
		"geo_diversity":   geoScore,
	}, nil
}
