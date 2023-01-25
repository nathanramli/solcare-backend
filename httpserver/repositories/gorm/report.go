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
