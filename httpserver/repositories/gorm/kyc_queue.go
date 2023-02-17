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

func (repo *kycQueueRepo) FindRecentKycRequest(ctx context.Context, address string) (*models.KycQueues, error) {
	kycQueue := new(models.KycQueues)
	err := repo.db.WithContext(ctx).Where("users_wallet_address = ?", address).Order("requested_at desc").Find(kycQueue).Error
	if err != nil {
		return nil, err
	}

	if kycQueue.Id == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return kycQueue, nil
}
