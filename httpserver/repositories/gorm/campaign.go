package gorm

import (
	"context"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
)

type campaignRepo struct {
	db *gorm.DB
}

func NewCampaignRepo(db *gorm.DB) repositories.CampaignRepo {
	return &campaignRepo{
		db: db,
	}
}

func (r *campaignRepo) SaveCampaign(ctx context.Context, campaign *models.Campaign) error {
	return r.db.WithContext(ctx).Save(campaign).Error
}

func (r *campaignRepo) FindCampaignById(ctx context.Context, address string) (*models.Campaign, error) {
	campaign := new(models.Campaign)
	err := r.db.WithContext(ctx).Where("address = ?", address).Take(campaign).Error
	return campaign, err
}

func (r *campaignRepo) FindCampaignByUser(ctx context.Context, userAddress string) ([]models.Campaign, error) {
	var campaigns []models.Campaign

	if err := r.db.WithContext(ctx).Where("owner_address = ?", userAddress).Find(&campaigns).Error; err != nil {
		return campaigns, err
	}
	return campaigns, nil
}
