package gorm

import (
	"context"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
)

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) repositories.CategoryRepo {
	return &categoryRepo{
		db: db,
	}
}

func (r *categoryRepo) FindAllCategories(ctx context.Context) ([]models.Categories, error) {
	var categories []models.Categories

	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}

func (r *categoryRepo) FindCategoryById(ctx context.Context, id uint) (*models.Categories, error) {
	category := new(models.Categories)
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(category).Error
	return category, err
}
