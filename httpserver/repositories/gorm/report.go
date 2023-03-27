package gorm

import (
	"context"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
)

type reportRepo struct {
	db *gorm.DB
}

func NewReportRepo(db *gorm.DB) repositories.ReportRepo {
	return &reportRepo{
		db: db,
	}
}

func (r *reportRepo) SaveReport(ctx context.Context, report *models.Reports) error {
	return r.db.WithContext(ctx).Save(report).Error
}

func (r *reportRepo) FindReportById(ctx context.Context, id uint) (*models.Reports, error) {
	report := new(models.Reports)
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(report).Error
	return report, err
}

func (r *reportRepo) FindGroupedReports(ctx context.Context) ([]map[string]interface{}, error) {
	type Result struct {
		CampaignAddress string
		Title           string
		Total           int
	}

	sliceMaps := make([]map[string]interface{}, 0)

	err := r.db.WithContext(ctx).Raw("SELECT rp.*, campaigns.title, campaigns.owner_address FROM (SELECT campaign_address, count(users_wallet_address) as total FROM reports GROUP BY campaign_address) rp JOIN campaigns ON campaigns.address = rp.campaign_address").Scan(&sliceMaps).Error

	return sliceMaps, err
}

func (r *reportRepo) FindReportsByAddress(ctx context.Context, address string) ([]models.Reports, error) {
	var reports []models.Reports

	if err := r.db.WithContext(ctx).Where("campaign_address = ?", address).Order("created_at desc").Find(&reports).Error; err != nil {
		return reports, err
	}
	return reports, nil
}
