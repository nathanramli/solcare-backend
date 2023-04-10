package views

type FindTransaction struct {
	Signature       string `json:"signature"`
	UserAddress     string `json:"userAddress"`
	CampaignAddress string `json:"campaignAddress"`
	CreatedAt       int64  `json:"createdAt"`
	Amount          uint64 `json:"amount"`
	Type            uint8  `json:"type"`
}
