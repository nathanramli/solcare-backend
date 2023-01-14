package services

import (
	"context"

	"github.com/nathanramli/solcare-backend/httpserver/controllers/params"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/views"
)

type UserSvc interface {
	Login(ctx context.Context, user *params.Login) *views.Response
}

type CampaignSvc interface {
	CreateCampaign(ctx context.Context, campaign *params.CreateCampaign) *views.Response
}

type CategoriesSvc interface {
	FindAllCategories(ctx context.Context) *views.Response
	FindCategoryById(ctx context.Context, categoryId uint) *views.Response
}