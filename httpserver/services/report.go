package services

import (
	"context"
	"github.com/nathanramli/solcare-backend/common"
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
	userRepo     repositories.UserRepo
}

func NewReportSvc(repo repositories.ReportRepo, campaignRepo repositories.CampaignRepo, userRepo repositories.UserRepo) ReportSvc {
	return &reportSvc{
		repo:         repo,
		campaignRepo: campaignRepo,
		userRepo:     userRepo,
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

func (svc *reportSvc) FindGroupedReports(ctx context.Context) *views.Response {
	reports, err := svc.repo.FindGroupedReports(ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	resp := make([]views.FindGroupedReports, len(reports))
	for i, v := range reports {
		resp[i].CampaignTitle = v["title"].(string)
		resp[i].CampaignAddress = v["campaign_address"].(string)
		resp[i].OwnerAddress = v["owner_address"].(string)
		resp[i].Total = v["total"].(int64)
	}
	return views.SuccessResponse(http.StatusOK, views.M_OK, resp)
}

func (svc *reportSvc) FindReportsByAddress(ctx context.Context, address string) *views.Response {
	reports, err := svc.repo.FindReportsByAddress(ctx, address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	resp := make([]views.FindReport, len(reports))
	for i, v := range reports {
		resp[i].Id = v.Id
		resp[i].Reporter = v.UsersWalletAddress
		resp[i].CampaignAddress = v.CampaignAddress
		resp[i].Description = v.Description
		resp[i].CreatedAt = v.CreatedAt.Unix()
	}
	return views.SuccessResponse(http.StatusOK, views.M_OK, resp)
}

func (svc *reportSvc) VerifyReport(ctx context.Context, params *params.VerifyReport) *views.Response {
	campaign, err := svc.campaignRepo.FindCampaignByAddress(ctx, params.CampaignAddress)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	user, err := svc.userRepo.FindUserByAddress(ctx, campaign.OwnerAddress)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	campaign.Delisted = common.GetBoolPointer(true)
	err = svc.campaignRepo.SaveCampaign(ctx, campaign)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	user.IsWarned = common.GetBoolPointer(true)
	err = svc.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, nil)
}
