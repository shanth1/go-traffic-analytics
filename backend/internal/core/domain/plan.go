package domain

// --- Plan Entity ---
type Plan struct {
	ID             string `json:"id"`               // e.g., "free", "pro"
	Name           string `json:"name"`             // "Starter", "Professional"
	PriceCents     int    `json:"price_cents"`      // 0, 990, 4990
	MaxLinks       int    `json:"max_links"`        // Лимит ссылок (-1 for infinity)
	MaxClicksMonth int    `json:"max_clicks_month"` // Лимит кликов
	CanExportData  bool   `json:"can_export_data"`  // Фича: экспорт CSV
	IsActive       bool   `json:"is_active"`
}
