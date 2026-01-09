package memoryrepo

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/consts"
)

type InMemoryAnalyticRepo struct {
	mu     sync.RWMutex
	events []*domain.ClickEvent
}

func NewAnalyticRepo() ports.AnalyticsRepository {
	return &InMemoryAnalyticRepo{
		events: make([]*domain.ClickEvent, 0),
	}
}

func (r *InMemoryAnalyticRepo) SaveBatch(_ context.Context, events []*domain.ClickEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, events...)
	return nil
}

func (r *InMemoryAnalyticRepo) filterClicks(filter domain.AnalyticsFilter) []*domain.ClickEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]*domain.ClickEvent, 0, len(r.events))
	for _, c := range r.events {
		if filter.LinkID != "" && c.LinkID != filter.LinkID {
			continue
		}
		if filter.CampaignID != "" && c.CampaignID != filter.CampaignID {
			continue
		}
		if !filter.From.IsZero() && c.Timestamp.Before(filter.From) {
			continue
		}
		if !filter.To.IsZero() && c.Timestamp.After(filter.To) {
			continue
		}
		filtered = append(filtered, c)
	}
	return filtered
}

func (r *InMemoryAnalyticRepo) CountTotal(_ context.Context, filter domain.AnalyticsFilter) (int64, error) {
	res := r.filterClicks(filter)
	return int64(len(res)), nil
}

// GetTimeSeriesGrouped needed for Streamgraph (Bucket by Time + Group by Dimension)
func (r *InMemoryAnalyticRepo) GetTimeSeriesGrouped(_ context.Context, filter domain.AnalyticsFilter, dimension string, interval time.Duration) ([]domain.StackedPoint, error) {
	events := r.filterClicks(filter)

	// Map: TimeBucket -> Category -> Count
	buckets := make(map[time.Time]map[string]int)

	for _, c := range events {
		bucketTime := c.Timestamp.Truncate(interval)

		if buckets[bucketTime] == nil {
			buckets[bucketTime] = make(map[string]int)
		}

		key := consts.Unknown
		switch dimension {
		case consts.OS:
			key = c.OS
		case consts.Browser:
			key = c.Browser
		case consts.Device:
			key = c.Device
		case consts.Country:
			key = c.Country
		}

		buckets[bucketTime][key]++
	}

	result := make([]domain.StackedPoint, 0, len(buckets))
	for t, values := range buckets {
		result = append(result, domain.StackedPoint{
			Time:   t,
			Values: values,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

type linkKey struct {
	srcID string
	tgtID string
}

// GetFlowData needed for Sankey (Flow Data)
//
// Transitions between stages:
// Referer -> OS (Layer 0 -> Layer 1)
// OS -> Country (Layer 1 -> Layer 2)
func (r *InMemoryAnalyticRepo) GetFlowData(_ context.Context, filter domain.AnalyticsFilter, stages []string) (*domain.SankeyData, error) {
	events := r.filterClicks(filter)

	nodesMap := make(map[string]domain.SankeyNode) // Key: "LayerIndex:Value"
	linksMap := make(map[linkKey]int)              // Key: "linkKey", Value: Count

	for _, c := range events {
		// Passable in stages in pairs
		for i := 0; i < len(stages)-1; i++ {
			sourceDim := stages[i]
			targetDim := stages[i+1]

			sourceVal := getDimensionValue(c, sourceDim)
			targetVal := getDimensionValue(c, targetDim)

			srcID := fmt.Sprintf("%d:%s", i, sourceVal) // "0:Google"
			tgtID := fmt.Sprintf("%d:%s", i+1, targetVal)

			// Nodes logic
			nodesMap[sourceVal] = domain.SankeyNode{ID: srcID, Layer: i}
			nodesMap[targetVal] = domain.SankeyNode{ID: tgtID, Layer: i + 1}

			// Links logic
			key := linkKey{srcID: srcID, tgtID: tgtID}
			linksMap[key]++
		}
	}

	// Formate Output
	data := &domain.SankeyData{
		Nodes: make([]domain.SankeyNode, 0, len(nodesMap)),
		Links: make([]domain.SankeyLink, 0, len(linksMap)),
	}

	for _, n := range nodesMap {
		data.Nodes = append(data.Nodes, n)
	}

	for key, count := range linksMap {
		data.Links = append(data.Links, domain.SankeyLink{
			Source: key.srcID,
			Target: key.tgtID,
			Value:  count,
		})
	}

	return data, nil
}

func (r *InMemoryAnalyticRepo) GetHeatmapData(_ context.Context, filter domain.AnalyticsFilter) ([]domain.HeatmapPoint, error) {
	events := r.filterClicks(filter)

	// Matrix [DayOfWeek][Hour]
	matrix := make(map[int]map[int]int)

	for _, c := range events {
		// time.Weekday: Sunday=0
		day := int(c.Timestamp.Weekday())
		hour := c.Timestamp.Hour()

		if matrix[day] == nil {
			matrix[day] = make(map[int]int)
		}
		matrix[day][hour]++
	}

	var result []domain.HeatmapPoint
	for d := 0; d < 7; d++ {
		for h := 0; h < 24; h++ {
			val := 0
			if matrix[d] != nil {
				val = matrix[d][h]
			}

			if val > 0 {
				result = append(result, domain.HeatmapPoint{
					DayOfWeek: d,
					Hour:      h,
					Count:     val,
				})
			}
		}
	}
	return result, nil
}

func (r *InMemoryAnalyticRepo) GetTopStats(_ context.Context, filter domain.AnalyticsFilter, dimension string, limit int) ([]domain.CategoryStat, error) {
	const op = "memoryrepo.InMemoryAnalyticRepo.GetTopStats"

	events := r.filterClicks(filter)
	if len(events) == 0 {
		return []domain.CategoryStat{}, nil
	}

	counts := make(map[string]int)

	const (
		ValueDirect  = "Direct"
		ValueUnknown = "Unknown"
	)

	for _, c := range events {
		var key string

		switch dimension {
		case consts.Browser:
			key = c.Browser
			if key == "" {
				key = ValueUnknown
			}
		case consts.OS:
			key = c.OS
			if key == "" {
				key = ValueUnknown
			}
		case consts.Device:
			key = c.Device
			if key == "" {
				key = ValueUnknown
			}
		case consts.Country:
			key = c.Country
			if key == "" {
				key = ValueUnknown
			}
		case consts.Referer:
			key = c.Referer
			if key == "" {
				key = ValueDirect
			}
		default:
			return nil, ops.E(op, ops.KindInvalid, fmt.Errorf("unsupported dimension for stats: %s", dimension))
		}

		counts[key]++
	}

	// Map -> Slice
	stats := make([]domain.CategoryStat, 0, len(counts))
	for k, v := range counts {
		stats = append(stats, domain.CategoryStat{Name: k, Value: v})
	}

	// Sort Descending
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Value > stats[j].Value
	})

	// Limit
	if limit > 0 && len(stats) > limit {
		result := make([]domain.CategoryStat, limit)
		copy(result, stats[:limit])

		topSum := 0
		for _, s := range result {
			topSum += s.Value
		}

		totalCount := 0
		for _, v := range counts {
			totalCount += v
		}

		othersCount := totalCount - topSum

		if othersCount > 0 {
			result = append(result, domain.CategoryStat{
				Name:  "Others",
				Value: othersCount,
			})
		}
		return result, nil
	}

	return stats, nil
}

func getDimensionValue(c *domain.ClickEvent, dim string) string {
	switch dim {
	case consts.Referer:
		return c.Referer
	case consts.OS:
		return c.OS
	case consts.Browser:
		return c.Browser
	case consts.Country:
		return c.Country
	case consts.Device:
		return c.Device
	}
	return "Other"
}
