package gorm

import (
	"context"
	"time"

	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repositories.UserRepo {
	return &userRepo{
		db: db,
	}
}

func (repo *userRepo) CreateUser(ctx context.Context, user *models.Users) error {
	user.CreatedAt = time.Now()
	return repo.db.WithContext(ctx).Create(user).Error
}

func (repo *userRepo) FindUserByAddress(ctx context.Context, address string) (*models.Users, error) {
	user := new(models.Users)
	return user, repo.db.WithContext(ctx).Where("wallet_address = ?", address).Take(user).Error
}

func (repo *userRepo) UpdateUser(ctx context.Context, user *models.Users) error {
	return repo.db.WithContext(ctx).Model(user).Updates(*user).Error
}
