package domain

import "time"

type LinkID string

type Link struct {
	ID         LinkID     `json:"id"`
	UserID     UserID     `json:"user_id"`
	CampaignID string     `json:"campaign_id"`
	Slug       string     `json:"slug"`
	TargetURL  string     `json:"target_url"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `json:"-"`
}
