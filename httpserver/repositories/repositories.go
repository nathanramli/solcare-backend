package repositories

import (
	"context"

	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.Users) error
	FindUserByAddress(ctx context.Context, address string) (*models.Users, error)
	FindUserById(ctx context.Context, id uint) (*models.Users, error)
	UpdateUser(ctx context.Context, user *models.Users) error
}

type CampaignRepo interface {
	SaveCampaign(ctx context.Context, campaign *models.Campaign) error
	FindCampaignById(ctx context.Context, address string) (*models.Campaign, error)
	FindCampaignByUser(ctx context.Context, userAddress string) ([]models.Campaign, error)
	FindAllCampaign(ctx context.Context, orders []string, limit int, offset int) ([]models.Campaign, error)
}

type CategoryRepo interface {
	FindAllCategories(ctx context.Context) ([]models.Categories, error)
	FindCategoryById(ctx context.Context, id uint) (*models.Categories, error)
}
