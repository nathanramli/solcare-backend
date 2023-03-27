package params

type CreateReport struct {
	CampaignAddress string `json:"campaignAddress" validate:"required"`
	Description     string `json:"description" validate:"required"`
}

type VerifyReport struct {
	CampaignAddress string `json:"campaignAddress" validate:"required"`
	IsAccepted      *bool  `json:"isAccepted"`
}
