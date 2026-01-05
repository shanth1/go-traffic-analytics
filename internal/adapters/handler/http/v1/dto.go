package v1

type CreateLinkRequest struct {
	CampaignID string `json:"campaign_id,omitempty"`
	TargetURL  string `json:"target_url" binding:"required" example:"https://google.com"`
}

type CreateCampaignRequest struct {
	Name string `json:"name" binding:"required"`
}
