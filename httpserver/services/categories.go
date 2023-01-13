package services

import (
	"context"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/views"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"gorm.io/gorm"
	"net/http"
)

type categorySvc struct {
	repo repositories.CategoryRepo
}

func NewCategorySvc(repo repositories.CategoryRepo) CategoriesSvc {
	return &categorySvc{
		repo: repo,
	}
}

func (svc *categorySvc) FindAllCategories(ctx context.Context) *views.Response {
	categories, err := svc.repo.FindAllCategories(ctx)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	resp := make([]views.FindAllCategories, len(categories))
	for i, category := range categories {
		r := views.FindAllCategories{
			Id:          category.Id,
			Name:        category.Name,
			Description: category.Description,
		}
		resp[i] = r
	}
	return views.SuccessResponse(http.StatusOK, views.M_OK, resp)
}

func (svc *categorySvc) FindCategoryById(ctx context.Context, categoryId uint) *views.Response {
	category, err := svc.repo.FindCategoryById(ctx, categoryId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, views.FindAllCategories{
		Id:          category.Id,
		Name:        category.Name,
		Description: category.Description,
	})
}
