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

type reportSvc struct {
	repo         repositories.ReportRepo
	campaignRepo repositories.CampaignRepo
}

func NewReportSvc(repo repositories.ReportRepo, campaignRepo repositories.CampaignRepo) ReportSvc {
	return &reportSvc{
		repo:         repo,
		campaignRepo: campaignRepo,
	}
}

func (svc *reportSvc) CreateReport(ctx context.Context, reporter string, params *params.CreateReport) *views.Response {
	_, err := svc.campaignRepo.FindCampaignByAddress(ctx, params.CampaignAddress)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	report := &models.Reports{
		CreatedAt:          time.Now(),
		UsersWalletAddress: reporter,
		CampaignAddress:    params.CampaignAddress,
		Description:        params.Description,
	}
	err = svc.repo.SaveReport(ctx, report)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, views.FindReport{
		Id:              report.Id,
		Reporter:        report.UsersWalletAddress,
		CampaignAddress: report.CampaignAddress,
		Description:     report.Description,
		CreatedAt:       report.CreatedAt.Unix(),
	})
}

func (svc *reportSvc) FindReportById(ctx context.Context, reportId uint) *views.Response {
	report, err := svc.repo.FindReportById(ctx, reportId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, views.FindReport{
		Id:              report.Id,
		Reporter:        report.UsersWalletAddress,
		CampaignAddress: report.CampaignAddress,
		Description:     report.Description,
		CreatedAt:       report.CreatedAt.Unix(),
	})
}
