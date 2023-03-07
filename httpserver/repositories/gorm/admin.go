package gorm

import (
	"context"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
)

type adminRepo struct {
	db *gorm.DB
}

func NewAdminRepo(db *gorm.DB) repositories.AdminRepo {
	return &adminRepo{
		db: db,
	}
}

func (r *adminRepo) FindAdminByAddress(ctx context.Context, address string) (*models.Admin, error) {
	admin := new(models.Admin)
	err := r.db.WithContext(ctx).Where("wallet_address = ?", address).Take(admin).Error
	return admin, err
}

func (r *adminRepo) FindAllAdmins(ctx context.Context) ([]models.Admin, error) {
	var admins []models.Admin

	if err := r.db.WithContext(ctx).Find(&admins).Error; err != nil {
		return admins, err
	}
	return admins, nil
}
