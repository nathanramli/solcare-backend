package gorm

import (
	"context"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
)

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) repositories.TransactionRepo {
	return &transactionRepo{
		db: db,
	}
}

func (r *transactionRepo) SaveTransaction(ctx context.Context, transaction *models.Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

func (r *transactionRepo) FindAllTransactionsByUser(ctx context.Context, address string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	if err := r.db.WithContext(ctx).Where("users_wallet_address = ?", address).Order("created_at desc").Find(&transactions).Error; err != nil {
		return transactions, err
	}
	return transactions, nil
}
