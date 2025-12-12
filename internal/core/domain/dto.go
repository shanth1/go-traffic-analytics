package domain

import "time"

// --- DTO для Аналитики (Visx Friendly) ---

// DataPoint - точка для линейных графиков
type TimeSeriesPoint struct {
	Time  time.Time `json:"time"`
	Value int       `json:"value"`
}

// StackedPoint - для Streamgraph (Time + категории)
// Пример: { Time: "12:00", "iOS": 10, "Android": 5 }
type StackedPoint struct {
	Time   time.Time      `json:"time"`
	Values map[string]int `json:"values"` // Динамические ключи (OS, Browser...)
}

// SankeyData - для диаграммы потоков
type SankeyData struct {
	Nodes []SankeyNode `json:"nodes"`
	Links []SankeyLink `json:"links"`
}

type SankeyNode struct {
	ID    string `json:"id"`    // Имя узла (напр. "USA" или "Mobile")
	Layer int    `json:"layer"` // Столбец (0 - Referer, 1 - Device...)
}

type SankeyLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}

// HeatmapPoint - для графика активности (День недели / Час)
// Visx Heatmap: x=Hour, y=Day, val=Count
type HeatmapPoint struct {
	DayOfWeek int `json:"day"`   // 0=Sun, 1=Mon, ..., 6=Sat
	Hour      int `json:"hour"`  // 0-23
	Count     int `json:"count"` // Интенсивность цвета
}

// CategoryStat - универсальная структура для Pie/Bar/Donut charts
// Например: [{Name: "Chrome", Value: 100}, {Name: "Firefox", Value: 50}]
type CategoryStat struct {
	Name  string  `json:"name"`
	Value int     `json:"value"`
	Share float64 `json:"share"` // Процент от общего (опционально, можно считать на фронте)
}

// RadarPoint - для Radar Chart (качество трафика)
type RadarPoint struct {
	Metric string `json:"metric"` // "Bot Score", "Unique IP", "Mobile %"
	Value  int    `json:"value"`  // 0-100
}

// --- Tree / Hierarchy (User -> Campaign -> Link) ---
type HierarchyNode struct {
	Name     string           `json:"name"`
	Type     string           `json:"type,omitempty"`  // "root", "campaign", "link"
	Value    int              `json:"value,omitempty"` // Например, клики (опционально)
	Children []*HierarchyNode `json:"children,omitempty"`
}
