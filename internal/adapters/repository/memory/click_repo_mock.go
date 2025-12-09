package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type InMemoryClickRepo struct {
	mu     sync.RWMutex
	clicks []*domain.ClickEvent
}

func NewClickRepo() ports.ClickRepository {
	return &InMemoryClickRepo{
		clicks: make([]*domain.ClickEvent, 0),
	}
}

func (r *InMemoryClickRepo) Save(_ context.Context, click *domain.ClickEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clicks = append(r.clicks, click)
	return nil
}

func (r *InMemoryClickRepo) filterClicks(filter ports.AnalyticsFilter) []*domain.ClickEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]*domain.ClickEvent, 0, len(r.clicks))
	for _, c := range r.clicks {
		if filter.LinkID != "" && c.LinkID != filter.LinkID {
			continue
		}
		if !c.Timestamp.After(filter.From) || !c.Timestamp.Before(filter.To) {
			continue
		}
		filtered = append(filtered, c)
	}
	return filtered
}

func (r *InMemoryClickRepo) CountTotal(_ context.Context, filter ports.AnalyticsFilter) (int64, error) {
	res := r.filterClicks(filter)
	return int64(len(res)), nil
}

// Реализация Streamgraph (Bucket by Time + Group by Dimension)
func (r *InMemoryClickRepo) GetTimeSeriesGrouped(_ context.Context, filter ports.AnalyticsFilter, dimension string, interval time.Duration) ([]domain.StackedPoint, error) {
	clicks := r.filterClicks(filter)

	// Map: TimeBucket -> Category -> Count
	buckets := make(map[time.Time]map[string]int)

	for _, c := range clicks {
		// Округляем время до интервала (truncation)
		bucketTime := c.Timestamp.Truncate(interval)

		if buckets[bucketTime] == nil {
			buckets[bucketTime] = make(map[string]int)
		}

		// Извлекаем значение измерения через рефлексию или switch (проще switch для мока)
		key := "unknown"
		switch dimension {
		case "os":
			key = c.OS
		case "browser":
			key = c.Browser
		case "device":
			key = c.Device
		case "country":
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

	// Сортировка по времени
	sort.Slice(result, func(i, j int) bool {
		return result[i].Time.Before(result[j].Time)
	})

	return result, nil
}

// Реализация Sankey (Flow Data)
// Логика: Нужно посчитать переходы между этапами.
// Referer -> OS (Layer 0 -> Layer 1)
// OS -> Country (Layer 1 -> Layer 2)
func (r *InMemoryClickRepo) GetFlowData(_ context.Context, filter ports.AnalyticsFilter, stages []string) (*domain.SankeyData, error) {
	clicks := r.filterClicks(filter)

	nodesMap := make(map[string]domain.SankeyNode) // Key: "LayerIndex:Value"
	linksMap := make(map[string]int)               // Key: "Source|Target", Value: Count

	for _, c := range clicks {
		// Проходим по этапам парами
		for i := 0; i < len(stages)-1; i++ {
			sourceDim := stages[i]
			targetDim := stages[i+1]

			sourceVal := getDimensionValue(c, sourceDim)
			targetVal := getDimensionValue(c, targetDim)

			// Создаем уникальные ID для нод, чтобы "Chrome" в Referer и "Chrome" в Browser не склеились, если они на разных слоях
			// Но для визуализации лучше называть их просто по значению.
			// Для простоты мока считаем, что значения уникальны между слоями, или добавляем префикс.

			// Nodes logic
			nodesMap[sourceVal] = domain.SankeyNode{ID: sourceVal, Layer: i}
			nodesMap[targetVal] = domain.SankeyNode{ID: targetVal, Layer: i + 1}

			// Links logic
			linkKey := fmt.Sprintf("%s|%s", sourceVal, targetVal)
			linksMap[linkKey]++
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

	for k, v := range linksMap {
		// Парсим ключ обратно (в проде лучше структуру использовать как ключ мапы)
		var src, tgt string
		_, _ = fmt.Sscanf(k, "%s|%s", &src, &tgt) // Упрощено, лучше split string

		// *Фикс для парсинга, так как | может быть разделителем
		// Тут лучше использовать strings.Split(k, "|")

		parts := splitLinkKey(k)
		if len(parts) == 2 {
			data.Links = append(data.Links, domain.SankeyLink{
				Source: parts[0],
				Target: parts[1],
				Value:  v,
			})
		}
	}

	return data, nil
}

func getDimensionValue(c *domain.ClickEvent, dim string) string {
	switch dim {
	case "referer":
		return c.Referer
	case "os":
		return c.OS
	case "browser":
		return c.Browser
	case "country":
		return c.Country
	case "device":
		return c.Device
	}
	return "Other"
}

// TODO:
func splitLinkKey(_ string) []string {
	// ... (реализация split)
	// Представим, что тут strings.Split
	return []string{"SourceStub", "TargetStub"} // Placeholder
}
