package domain

import "time"

type CreateLinkCmd struct {
	UserID     UserID
	CampaignID string
	TargetURL  string
	CustomSlug string
}

type LinkFilter struct {
	UserID     UserID
	CampaignID string
	Search     string // slug / target_url
	IsActive   *bool
	Limit      int
	Offset     int
}

// AnalyticsFilter для всех отчетов
type AnalyticsFilter struct {
	UserID     UserID
	LinkID     LinkID
	CampaignID string
	From       time.Time
	To         time.Time
}

// RequestMetadata — контекст клика (передается из HTTP хендлера)
type RequestMetadata struct {
	IP        string
	UserAgent string
	Referer   string
}

type GeoLocation struct {
	CountryCode string
	Country     string
	City        string
	Latitude    float64
	Longitude   float64
}
