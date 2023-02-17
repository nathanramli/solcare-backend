package repositories

import (
	"context"

	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.Users) error
	FindUserByAddress(ctx context.Context, address string) (*models.Users, error)
	UpdateUser(ctx context.Context, user *models.Users) error
}

type CampaignRepo interface {
	SaveCampaign(ctx context.Context, campaign *models.Campaign) error
	FindCampaignByAddress(ctx context.Context, address string) (*models.Campaign, error)
	FindCampaignByUser(ctx context.Context, userAddress string) ([]models.Campaign, error)
	FindAllCampaign(ctx context.Context, orders []string, limit int, offset int, filters string) ([]models.Campaign, error)
}

type ProposalRepo interface {
	SaveProposal(ctx context.Context, proposal *models.Proposal) error
	FindProposalByAddress(ctx context.Context, address string) (*models.Proposal, error)
}

type CategoryRepo interface {
	FindAllCategories(ctx context.Context) ([]models.Categories, error)
	FindCategoryById(ctx context.Context, id uint) (*models.Categories, error)
}

type ReportRepo interface {
	SaveReport(ctx context.Context, report *models.Reports) error
	FindReportById(ctx context.Context, id uint) (*models.Reports, error)
}

type KycQueueRepo interface {
	SaveKycQueue(ctx context.Context, request *models.KycQueues) error
	FindRecentKycRequest(ctx context.Context, address string) (*models.KycQueues, error)
}
