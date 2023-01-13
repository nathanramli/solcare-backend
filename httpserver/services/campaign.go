package services

import (
	"context"
	"errors"
	"github.com/gagliardetto/solana-go"
	"github.com/gin-gonic/gin"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/params"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/views"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type campaignSvc struct {
	repo     repositories.CampaignRepo
	cateRepo repositories.CategoryRepo
}

func NewCampaignSvc(repo repositories.CampaignRepo, cateRepo repositories.CategoryRepo) CampaignSvc {
	return &campaignSvc{
		repo:     repo,
		cateRepo: cateRepo,
	}
}

func (svc *campaignSvc) CreateCampaign(ctx context.Context, params *params.CreateCampaign) *views.Response {
	pubkey, err := solana.PublicKeyFromBase58(params.Address)
	if err != nil {
		return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
	}

	_, err = solana.PublicKeyFromBase58(params.OwnerAddress)
	if err != nil {
		return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
	}

	_, err = svc.cateRepo.FindCategoryById(ctx, params.CategoryId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	_, err = svc.repo.FindCampaignById(ctx, params.Address)
	if err != nil && err != gorm.ErrRecordNotFound {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	} else if err == nil {
		return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, errors.New("campaign already exist"))
	}

	err = svc.repo.CreateCampaign(ctx, &models.Campaign{
		Title:        params.Title,
		Description:  params.Description,
		Address:      params.Address,
		OwnerAddress: params.OwnerAddress,
		CategoryId:   params.CategoryId,
	})
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	fileNameSplits := strings.Split(params.Banner.Filename, ".")
	ext := fileNameSplits[len(fileNameSplits)-1]

	err = ctx.(*gin.Context).SaveUploadedFile(params.Banner, "./resources/"+pubkey.String()+"."+ext)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, nil)
}
