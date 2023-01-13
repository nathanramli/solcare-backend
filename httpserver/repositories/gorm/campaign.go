package gorm

import (
	"context"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
	"time"
)

type campaignRepo struct {
	db *gorm.DB
}

func NewCampaignRepo(db *gorm.DB) repositories.CampaignRepo {
	return &campaignRepo{
		db: db,
	}
}

func (r *campaignRepo) CreateCampaign(ctx context.Context, campaign *models.Campaign) error {
	campaign.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(campaign).Error
}

func (r *campaignRepo) FindCampaignById(ctx context.Context, address string) (*models.Campaign, error) {
	campaign := new(models.Campaign)
	err := r.db.WithContext(ctx).Where("address = ?", address).Take(campaign).Error
	return campaign, err
}
