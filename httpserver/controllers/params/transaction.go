package params

type CreateTransaction struct {
	Signature       string `json:"signature" validate:"required"`
	CampaignAddress string `json:"campaignAddress" validate:"required"`
	Amount          uint64 `json:"amount"`
	Type            uint8  `json:"type"`
}
