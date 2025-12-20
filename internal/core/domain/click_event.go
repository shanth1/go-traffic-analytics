package domain

import "time"

type ClickEvent struct {
	ID         string    `json:"id"`
	LinkID     string    `json:"link_id"`
	CampaignID string    `json:"campaign_id"`
	Timestamp  time.Time `json:"timestamp"`
	IP         string    `json:"ip"`
	Country    string    `json:"country"` // ISO (US, DE, RU)
	City       string    `json:"city"`
	OS         string    `json:"os"`      // iOS, Android, Windows
	Browser    string    `json:"browser"` // Chrome, Safari
	Device     string    `json:"device"`  // Mobile, Desktop
	Referer    string    `json:"referer"` // Google, Facebook, Direct
}
