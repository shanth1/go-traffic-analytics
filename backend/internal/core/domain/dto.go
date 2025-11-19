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
