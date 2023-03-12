package params

import "mime/multipart"

type CreateCampaign struct {
	Address      string                `form:"address" validate:"required"`
	OwnerAddress string                `form:"ownerAddress" validate:"required"`
	Title        string                `form:"title" validate:"required"`
	Description  string                `form:"description" validate:"required"`
	CategoryId   uint                  `form:"categoryId" validate:"required"`
	Banner       *multipart.FileHeader `form:"banner" validate:"required"`
}

type UploadEvidence struct {
	CampaignAddress string                `form:"campaignAddress" validate:"required"`
	Attachment      *multipart.FileHeader `form:"attachment" validate:"required"`
}

type VerifyEvidence struct {
	Address    string `json:"address" validate:"required"`
	IsApproved *bool  `json:"isApproved"`
}
