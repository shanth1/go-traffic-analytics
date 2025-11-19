package domain

import "time"

type Link struct {
	ID        string    `json:"id"`
	Slug      string    `json:"slug"`
	TargetURL string    `json:"target_url"`
	CreatedAt time.Time `json:"created_at"`
}
