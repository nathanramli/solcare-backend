package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nathanramli/solcare-backend/common"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/params"
	"github.com/nathanramli/solcare-backend/httpserver/services"
	"net/http"
	"strconv"
)

type ReportController struct {
	svc services.ReportSvc
}

func NewReportController(svc services.ReportSvc) *ReportController {
	return &ReportController{
		svc: svc,
	}
}

func (control *ReportController) CreateReport(ctx *gin.Context) {
	var req params.CreateReport
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = validator.New().Struct(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	claims, _ := ctx.Get("userData")
	userData := claims.(*common.CustomClaims)

	response := control.svc.CreateReport(ctx, userData.Address, &req)
	WriteJsonResponse(ctx, response)
}

func (control *ReportController) FindReportById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := control.svc.FindReportById(ctx, uint(id))
	WriteJsonResponse(ctx, response)
}

func (control *ReportController) FindGroupedReports(ctx *gin.Context) {
	response := control.svc.FindGroupedReports(ctx)
	WriteJsonResponse(ctx, response)
}

func (control *ReportController) FindReportsByAddress(ctx *gin.Context) {
	response := control.svc.FindReportsByAddress(ctx, ctx.Param("address"))
	WriteJsonResponse(ctx, response)
}

func (control *ReportController) VerifyReport(ctx *gin.Context) {
	var req params.VerifyReport
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = validator.New().Struct(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := control.svc.VerifyReport(ctx, &req)
	WriteJsonResponse(ctx, response)
}
