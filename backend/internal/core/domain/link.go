package domain

import "time"

type Link struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CampaignID string    `json:"campaign_id"`
	Slug       string    `json:"slug"`
	TargetURL  string    `json:"target_url"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}
