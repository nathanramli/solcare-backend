package services

import (
	"context"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/params"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/views"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type transactionSvc struct {
	repo         repositories.TransactionRepo
	campaignRepo repositories.CampaignRepo
	userRepo     repositories.UserRepo
}

func NewTransactionSvc(repo repositories.TransactionRepo, campaignRepo repositories.CampaignRepo, userRepo repositories.UserRepo) TransactionSvc {
	return &transactionSvc{
		repo:         repo,
		campaignRepo: campaignRepo,
		userRepo:     userRepo,
	}
}

func (svc *transactionSvc) CreateTransaction(ctx context.Context, userAddress string, params *params.CreateTransaction) *views.Response {
	_, err := svc.campaignRepo.FindCampaignByAddress(ctx, params.CampaignAddress)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	transaction := &models.Transaction{
		CreatedAt:          time.Now(),
		UsersWalletAddress: userAddress,
		CampaignAddress:    params.CampaignAddress,
		Signature:          params.Signature,
		Amount:             params.Amount,
		Type:               params.Type,
	}
	err = svc.repo.SaveTransaction(ctx, transaction)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, views.FindTransaction{
		Signature:       transaction.Signature,
		UserAddress:     transaction.UsersWalletAddress,
		CampaignAddress: transaction.CampaignAddress,
		CreatedAt:       transaction.CreatedAt.Unix(),
		Amount:          transaction.Amount,
		Type:            transaction.Type,
	})
}

func (svc *transactionSvc) FindAllTransactionsByUser(ctx context.Context, address string) *views.Response {
	transactions, err := svc.repo.FindAllTransactionsByUser(ctx, address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	resp := make([]views.FindTransaction, len(transactions))
	for i, v := range transactions {
		resp[i] = views.FindTransaction{
			Signature:       v.Signature,
			UserAddress:     v.UsersWalletAddress,
			CampaignAddress: v.CampaignAddress,
			CreatedAt:       v.CreatedAt.Unix(),
			Amount:          v.Amount,
			Type:            v.Type,
		}
	}
	return views.SuccessResponse(http.StatusOK, views.M_OK, resp)
}
