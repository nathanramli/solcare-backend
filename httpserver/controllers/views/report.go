package views

type FindReport struct {
	Id              uint   `json:"id"`
	Reporter        string `json:"reporter"`
	CampaignAddress string `json:"campaignAddress"`
	Description     string `json:"description"`
	CreatedAt       int64  `json:"createdAt"`
}
