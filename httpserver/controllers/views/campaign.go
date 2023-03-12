package views

type FindAllCampaigns struct {
	Address      string `json:"address"`
	CreatedAt    int64  `json:"createdAt"`
	OwnerAddress string `json:"ownerAddress"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CategoryId   uint   `json:"categoryId"`
	Status       uint8  `json:"status"`
	Banner       string `json:"banner"`
	Delisted     bool   `json:"delisted"`
}

type FindAllCampaignsWithEvidence struct {
	Address      string `json:"address"`
	OwnerAddress string `json:"ownerAddress"`
	Status       uint8  `json:"status"`
	Evidence     string `json:"evidence"`
	Delisted     bool   `json:"delisted"`
}
