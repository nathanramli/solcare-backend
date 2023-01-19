package params

import "mime/multipart"

type CreateProposal struct {
	Address         string                `form:"address" validate:"required"`
	CampaignAddress string                `form:"campaignAddress" validate:"required"`
	Attachment      *multipart.FileHeader `form:"attachment" validate:"required"`
}
