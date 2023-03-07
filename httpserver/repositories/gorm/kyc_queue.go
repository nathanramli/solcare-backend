package gorm

import (
	"context"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
)

type kycQueueRepo struct {
	db *gorm.DB
}

func NewKyqQueueRepo(db *gorm.DB) repositories.KycQueueRepo {
	return &kycQueueRepo{
		db: db,
	}
}

func (repo *kycQueueRepo) SaveKycQueue(ctx context.Context, request *models.KycQueues) error {
	return repo.db.WithContext(ctx).Save(request).Error
}

func (repo *kycQueueRepo) FindKycRequestByUser(ctx context.Context, address string) (*models.KycQueues, error) {
	kycQueue := new(models.KycQueues)
	err := repo.db.WithContext(ctx).Preload("Users").Where("users_wallet_address = ?", address).Take(kycQueue).Error
	return kycQueue, err
}

func (repo *kycQueueRepo) FindAllKycRequest(ctx context.Context, status int) ([]models.KycQueues, error) {
	var kycQueues []models.KycQueues

	if err := repo.db.WithContext(ctx).Preload("Users").Where("status = ?", status).Find(&kycQueues).Error; err != nil {
		return kycQueues, err
	}
	return kycQueues, nil
}
