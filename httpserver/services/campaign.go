package services

import (
	"context"
	"errors"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gin-gonic/gin"
	"github.com/nathanramli/solcare-backend/config"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/params"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/views"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
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

func (svc *campaignSvc) FindCampaignByUser(ctx context.Context, userAddress string) *views.Response {
	campaigns, err := svc.repo.FindCampaignByUser(ctx, userAddress)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	resp := make([]views.FindAllCampaigns, len(campaigns))
	for i, campaign := range campaigns {
		r := views.FindAllCampaigns{
			Address:      campaign.Address,
			CreatedAt:    campaign.CreatedAt.Unix(),
			OwnerAddress: campaign.OwnerAddress,
			Title:        campaign.Title,
			Description:  campaign.Description,
			CategoryId:   campaign.CategoryId,
			Status:       campaign.Status,
			Banner:       "resources/" + campaign.Banner,
			Delisted:     campaign.Delisted,
		}
		resp[i] = r
	}
	return views.SuccessResponse(http.StatusOK, views.M_OK, resp)
}

func (svc *campaignSvc) FindCampaignByAddress(ctx context.Context, address string) *views.Response {
	campaign, err := svc.repo.FindCampaignByAddress(ctx, address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, views.FindAllCampaigns{
		Address:      campaign.Address,
		CreatedAt:    campaign.CreatedAt.Unix(),
		OwnerAddress: campaign.OwnerAddress,
		Title:        campaign.Title,
		Description:  campaign.Description,
		CategoryId:   campaign.CategoryId,
		Status:       campaign.Status,
		Banner:       "resources/" + campaign.Banner,
		Delisted:     campaign.Delisted,
	})
}

func (svc *campaignSvc) FindAllCampaign(ctx context.Context, order string, categoryId int, search string, offset int) *views.Response {
	orders := make([]string, 0)
	if order == "newest" {
		orders = append(orders, "created_at desc")
	} else if order == "oldest" {
		orders = append(orders, "created_at asc")
	} else {
		// default
		orders = append(orders, "created_at desc")
	}

	filters := ""
	if categoryId != 0 {
		filters += "category_id = '" + strconv.Itoa(categoryId) + "'"
	}

	if search != "" {
		if filters != "" {
			filters += " AND "
		}
		filters += "LOWER(title) LIKE '%" + strings.ToLower(search) + "%'"
	}

	campaigns, err := svc.repo.FindAllCampaign(ctx, orders, 20, offset, filters)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	resp := make([]views.FindAllCampaigns, len(campaigns))
	for i, campaign := range campaigns {
		r := views.FindAllCampaigns{
			Address:      campaign.Address,
			CreatedAt:    campaign.CreatedAt.Unix(),
			OwnerAddress: campaign.OwnerAddress,
			Title:        campaign.Title,
			Description:  campaign.Description,
			CategoryId:   campaign.CategoryId,
			Status:       campaign.Status,
			Banner:       "resources/" + campaign.Banner,
			Delisted:     campaign.Delisted,
		}
		resp[i] = r
	}
	return views.SuccessResponse(http.StatusOK, views.M_OK, resp)
}

func (svc *campaignSvc) CreateCampaign(ctx context.Context, params *params.CreateCampaign) *views.Response {
	pubkey, err := solana.PublicKeyFromBase58(params.Address)
	if err != nil {
		return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
	}

	acc, err := config.RpcClient.GetAccountInfo(ctx, pubkey)
	if err != nil && err != rpc.ErrNotFound {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	if acc != nil {
		return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, errors.New("campaign already exist"))
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

	fileNameSplits := strings.Split(params.Banner.Filename, ".")
	ext := fileNameSplits[len(fileNameSplits)-1]

	err = ctx.(*gin.Context).SaveUploadedFile(params.Banner, "./resources/"+pubkey.String()+"."+ext)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	err = svc.repo.SaveCampaign(ctx, &models.Campaign{
		Title:        params.Title,
		Description:  params.Description,
		Address:      params.Address,
		OwnerAddress: params.OwnerAddress,
		CategoryId:   params.CategoryId,
		Banner:       pubkey.String() + "." + ext,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, nil)
}
