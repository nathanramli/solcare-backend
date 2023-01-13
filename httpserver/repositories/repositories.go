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
	CreateCampaign(ctx context.Context, campaign *models.Campaign) error
	FindCampaignById(ctx context.Context, address string) (*models.Campaign, error)
}

type CategoryRepo interface {
	FindAllCategories(ctx context.Context) ([]models.Categories, error)
	FindCategoryById(ctx context.Context, id uint) (*models.Categories, error)
}
