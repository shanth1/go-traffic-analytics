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
	// TODO: Unique Users, Top Country etc.
	return map[string]interface{}{
		"total_clicks": total,
	}, nil
}

func (s *AnalyticsService) GetStreamGraphData(ctx context.Context, filter ports.AnalyticsFilter, groupBy string) ([]domain.StackedPoint, error) {
	interval := time.Hour * 24 // Day by default
	return s.clickRepo.GetTimeSeriesGrouped(ctx, filter, groupBy, interval)
}

func (s *AnalyticsService) GetSankeyData(ctx context.Context, filter ports.AnalyticsFilter, stages []string) (*domain.SankeyData, error) {
	return s.clickRepo.GetFlowData(ctx, filter, stages)
}

func (s *AnalyticsService) GetGeoDistribution(ctx context.Context, filter ports.AnalyticsFilter) (interface{}, error) {
	data, err := s.clickRepo.GetTimeSeriesGrouped(ctx, filter, "country", time.Hour*24*365*10) // Hack for total sum
	if err != nil {
		return nil, err
	}

	totals := make(map[string]int)
	for _, point := range data {
		for country, count := range point.Values {
			totals[country] += count
		}
	}
	return totals, nil
}

// TODO:
func (s *AnalyticsService) GetTrafficQuality(ctx context.Context, filter ports.AnalyticsFilter) (map[string]int, error) {
	return map[string]int{
		"high_quality": 85,
		"suspicious":   10,
		"bot":          5,
	}, nil
}
