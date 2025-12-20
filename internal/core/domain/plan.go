package domain

// --- Plan Entity ---
type Plan struct {
	ID             string `json:"id"`   // e.g., "free", "pro"
	Name           string `json:"name"` // "Starter", "Professional"
	Description    string `json:"description"`
	PriceCents     int    `json:"price_cents"` // 0, 990, 4990
	MaxLinks       int    `json:"max_links"`   // -1 for infinity
	MaxClicksMonth int    `json:"max_clicks_month"`
	CanExportData  bool   `json:"can_export_data"`
	IsActive       bool   `json:"is_active"`
}
