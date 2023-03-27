package views

type FindReport struct {
	Id              uint   `json:"id"`
	Reporter        string `json:"reporter"`
	CampaignAddress string `json:"campaignAddress"`
	Description     string `json:"description"`
	CreatedAt       int64  `json:"createdAt"`
}

type FindGroupedReports struct {
	CampaignAddress string `json:"campaignAddress"`
	OwnerAddress    string `json:"ownerAddress"`
	CampaignTitle   string `json:"campaignTitle"`
	Total           int64  `json:"total"`
}
