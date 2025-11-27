package v1

type CreateLinkRequest struct {
	CampaignID string `json:"campaign_id"`
	TargetURL  string `json:"target_url"`
	CustomSlug string `json:"custom_slug"`
}

type CreateCampaignRequest struct {
	Name string `json:"name"`
}
