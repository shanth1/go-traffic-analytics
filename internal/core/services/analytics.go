package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/consts"
	"github.com/xuri/excelize/v2"
)

type AnalyticsService struct {
	repo     ports.AnalyticsRepository
	userRepo ports.UserRepository
	planRepo ports.PlanRepository
}

func NewAnalyticsService(analyticsRepo ports.AnalyticsRepository, userRepo ports.UserRepository, planRepo ports.PlanRepository) *AnalyticsService {
	return &AnalyticsService{
		repo:     analyticsRepo,
		userRepo: userRepo,
		planRepo: planRepo,
	}
}

func (s *AnalyticsService) GetSummary(ctx context.Context, filter domain.AnalyticsFilter) (*domain.Summary, error) {
	total, err := s.repo.CountTotal(ctx, filter)
	if err != nil {
		return nil, err
	}

	browsers, err := s.repo.GetTopStats(ctx, filter, consts.Browser, 5)
	if err != nil {
		// TODO: logging
		browsers = []domain.CategoryStat{}
	}

	osStats, err := s.repo.GetTopStats(ctx, filter, consts.OS, 5)
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

func (s *AnalyticsService) GetStreamGraphData(ctx context.Context, filter domain.AnalyticsFilter, groupBy string) ([]domain.StackedPoint, error) {
	duration := filter.To.Sub(filter.From)
	interval := time.Hour * 24 // Day by default

	if duration < time.Hour*48 {
		interval = time.Hour
	}

	return s.repo.GetTimeSeriesGrouped(ctx, filter, groupBy, interval)
}

func (s *AnalyticsService) GetSankeyData(ctx context.Context, filter domain.AnalyticsFilter, stages []string) (*domain.SankeyData, error) {
	if len(stages) == 0 {
		stages = []string{consts.Referer, consts.Device, consts.Country}
	}

	return s.repo.GetFlowData(ctx, filter, stages)
}

func (s *AnalyticsService) GetGeoDistribution(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.GeoStat, error) {
	stats, err := s.repo.GetTopStats(ctx, filter, consts.Country, 200)
	if err != nil {
		return nil, err
	}

	result := make([]domain.GeoStat, 0, len(stats))
	for _, stat := range stats {
		result = append(result, domain.GeoStat{
			Country: stat.Name,
			Value:   stat.Value,
		})
	}

	return result, nil
}

func (s *AnalyticsService) GetHeatmapData(ctx context.Context, filter domain.AnalyticsFilter) ([]domain.HeatmapPoint, error) {
	return s.repo.GetHeatmapData(ctx, filter)
}

func (s *AnalyticsService) GetCategoryStats(ctx context.Context, filter domain.AnalyticsFilter, dimension string) ([]domain.CategoryStat, error) {
	return s.repo.GetTopStats(ctx, filter, dimension, 10)
}

func (s *AnalyticsService) GetTrafficQuality(ctx context.Context, filter domain.AnalyticsFilter) (*domain.TrafficQuality, error) {
	total, err := s.repo.CountTotal(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("counting total clicks: %w", err)
	}

	if total == 0 {
		return &domain.TrafficQuality{
			HumanScore: 100,
		}, nil
	}

	// 1. Mobile Friendly
	devices, err := s.repo.GetTopStats(ctx, filter, consts.Device, 100)
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
	osStats, err := s.repo.GetTopStats(ctx, filter, consts.OS, 100)
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
	browserStats, err := s.repo.GetTopStats(ctx, filter, consts.Browser, 100)
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
	countries, err := s.repo.GetTopStats(ctx, filter, consts.Country, 100)
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

func (s *AnalyticsService) ExportData(ctx context.Context, filter domain.AnalyticsFilter) (io.Reader, string, error) {
	user, err := s.userRepo.FindByID(ctx, filter.UserID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user for export check: %w", err)
	}

	plan, err := s.planRepo.FindByID(ctx, user.PlanID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get plan: %w", err)
	}

	if !plan.CanExportData {
		return nil, "", fmt.Errorf("export feature is not available on plan %s", plan.Name)
	}

	events, err := s.repo.GetRawEvents(ctx, filter)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch raw events: %w", err)
	}

	f := excelize.NewFile()
	sheetName := "Clicks Data"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, "", err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{"Timestamp", "IP", "Country", "City", "OS", "Browser", "Device", "Referer", "Campaign ID", "Link ID"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	f.SetRowStyle(sheetName, 1, 1, headerStyle)

	for i, event := range events {
		row := i + 2

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), event.Timestamp.Format(time.RFC3339))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), event.IP)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), event.Country)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), event.City)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), event.OS)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), event.Browser)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), event.Device)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), event.Referer)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), event.CampaignID)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), string(event.LinkID))
	}

	var b bytes.Buffer
	if err := f.Write(&b); err != nil {
		return nil, "", fmt.Errorf("failed to write excel buffer: %w", err)
	}

	fileName := fmt.Sprintf("analytics_export_%s_%s.xlsx",
		filter.From.Format("20060102"),
		filter.To.Format("20060102"),
	)

	return &b, fileName, nil
}
