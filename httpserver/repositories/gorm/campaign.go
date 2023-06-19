package gorm

import (
	"context"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
	"strings"
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

func (r *campaignRepo) FindCampaignByAddress(ctx context.Context, address string) (*models.Campaign, error) {
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

func (r *campaignRepo) FindAllCampaign(ctx context.Context, orders []string, limit int, offset int, filters string) ([]models.Campaign, error) {
	var campaigns []models.Campaign

	query := r.db.WithContext(ctx).Limit(limit).Offset(offset)

	if err := query.Order(strings.Join(orders, ", ")).Find(&campaigns, filters).Error; err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *campaignRepo) FindAllCampaignWithEvidence(ctx context.Context) ([]models.Campaign, error) {
	var campaigns []models.Campaign

	if err := r.db.WithContext(ctx).Where("status = ?", models.EVIDENCE_STATUS_REQUESTED).Find(&campaigns).Error; err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *campaignRepo) CountTotalCampaigns(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&models.Campaign{}).Count(&total).Error
	return total, err
}

func (r *campaignRepo) CountTotalSuccessCampaigns(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&models.Campaign{}).Where("status = ?", models.EVIDENCE_STATUS_SUCCESS).Count(&total).Error
	return total, err
}

func (r *campaignRepo) CountTotalFailedCampaigns(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&models.Campaign{}).Where("status = ?", models.EVIDENCE_STATUS_FAILED).Count(&total).Error
	return total, err
}

func (r *campaignRepo) CountTotalDelistedCampaigns(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&models.Campaign{}).Where("delisted = true").Count(&total).Error
	return total, err
}
