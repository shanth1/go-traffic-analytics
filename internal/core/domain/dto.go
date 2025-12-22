package domain

import "time"

// --- Analytics DTOs (Visx Friendly) ---

// TimeSeriesPoint - point for line charts
type TimeSeriesPoint struct {
	Time  time.Time `json:"time"`
	Value int       `json:"value"`
}

// StackedPoint - for Streamgraph (Time + categories)
// Example: { Time: "12:00", "iOS": 10, "Android": 5 }
type StackedPoint struct {
	Time   time.Time      `json:"time"`
	Values map[string]int `json:"values"` // Dynamic keys (OS, Browser...)
}

// SankeyData - for flow diagrams
type SankeyData struct {
	Nodes []SankeyNode `json:"nodes"`
	Links []SankeyLink `json:"links"`
}

type SankeyNode struct {
	ID    string `json:"id"`    // Node name (e.g., "USA" or "Mobile")
	Layer int    `json:"layer"` // Column (0 - Referer, 1 - Device...)
}

type SankeyLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}

// HeatmapPoint - for activity chart (Day of Week / Hour)
// Visx Heatmap: x=Hour, y=Day, val=Count
type HeatmapPoint struct {
	DayOfWeek int `json:"day"`   // 0=Sun, 1=Mon, ..., 6=Sat
	Hour      int `json:"hour"`  // 0-23
	Count     int `json:"count"` // Color intensity
}

// CategoryStat - universal structure for Pie/Bar/Donut charts
// Example: [{Name: "Chrome", Value: 100}, {Name: "Firefox", Value: 50}]
type CategoryStat struct {
	Name  string  `json:"name"`
	Value int     `json:"value"`
	Share float64 `json:"share"` // Percentage of total (optional, can be calculated on frontend)
}

// RadarPoint - for Radar Chart (traffic quality)
type RadarPoint struct {
	Metric string `json:"metric"` // "Bot Score", "Unique IP", "Mobile %"
	Value  int    `json:"value"`  // 0-100
}

// --- Tree / Hierarchy (User -> Campaign -> Link) ---
type HierarchyNode struct {
	Name     string           `json:"name"`
	Type     string           `json:"type,omitempty"`  // "root", "campaign", "link"
	Value    int              `json:"value,omitempty"` // E.g., clicks (optional)
	Children []*HierarchyNode `json:"children,omitempty"`
}

// --- Traffic Quality ---
type TrafficQuality struct {
	MobileFriendlyScore int  `json:"mobile_friendly_score"`
	BotScore            int  `json:"bot_score"`
	HumanScore          int  `json:"human_score"`
	GeoDiversityScore   int  `json:"geo_diversity_score"`
	IsSuspicious        bool `json:"is_suspicious"`
}

type Summary struct {
	TotalClicks int64          `json:"total_clicks"`
	TopBrowsers []CategoryStat `json:"top_browsers"`
	TopOS       []CategoryStat `json:"top_os"`
}

type GeoStat struct {
	Country string `json:"country"`
	Value   int    `json:"value"`
}
