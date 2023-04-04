package services

import (
	"context"

	"github.com/nathanramli/solcare-backend/httpserver/controllers/params"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/views"
)

type UserSvc interface {
	Login(ctx context.Context, user *params.Login) *views.Response
	UpdateUser(ctx context.Context, address string, params *params.UpdateUser) *views.Response
	UpdateAvatar(ctx context.Context, address string, params *params.UpdateUserAvatar) *views.Response
	RequestKyc(ctx context.Context, address string, params *params.RequestKyc) *views.Response
	VerifyKyc(ctx context.Context, params *params.VerifyKyc) *views.Response
	RemoveKyc(ctx context.Context, address string) *views.Response
	FindUserByAddress(ctx context.Context, address string) *views.Response
	FindAllUsers(ctx context.Context) *views.Response
	FindKycRequestByUser(ctx context.Context, address string) *views.Response
	FindAllKycRequest(ctx context.Context, status int) *views.Response
}

type CampaignSvc interface {
	CreateCampaign(ctx context.Context, campaign *params.CreateCampaign) *views.Response
	FindCampaignByUser(ctx context.Context, userAddress string) *views.Response
	FindCampaignByAddress(ctx context.Context, address string) *views.Response
	FindAllCampaign(ctx context.Context, order string, category int, search string, offset int) *views.Response
	FindAllCampaignWithEvidence(ctx context.Context) *views.Response
	FetchCampaignSummary(ctx context.Context) *views.Response

	CreateProposal(ctx context.Context, campaign *params.CreateProposal) *views.Response
	FindProposalByAddress(ctx context.Context, address string) *views.Response

	UploadEvidence(ctx context.Context, params *params.UploadEvidence) *views.Response
	VerifyEvidence(ctx context.Context, params *params.VerifyEvidence) *views.Response
}

type CategoriesSvc interface {
	FindAllCategories(ctx context.Context) *views.Response
	FindCategoryById(ctx context.Context, categoryId uint) *views.Response
}

type ReportSvc interface {
	CreateReport(ctx context.Context, reporter string, params *params.CreateReport) *views.Response
	FindReportById(ctx context.Context, id uint) *views.Response
	FindGroupedReports(ctx context.Context) *views.Response
	FindReportsByAddress(ctx context.Context, address string) *views.Response
	VerifyReport(ctx context.Context, params *params.VerifyReport) *views.Response
}
