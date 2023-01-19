package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/params"
	"github.com/nathanramli/solcare-backend/httpserver/services"
	"net/http"
	"strconv"
)

type CampaignController struct {
	svc services.CampaignSvc
}

func NewCampaignController(svc services.CampaignSvc) *CampaignController {
	return &CampaignController{
		svc: svc,
	}
}

func (control *CampaignController) FindAllCampaign(ctx *gin.Context) {
	var err error
	var (
		offset     = 0
		categoryId = 0
	)

	offsetPar, exist := ctx.GetQuery("offset")
	if exist && offsetPar != "" {
		offset, err = strconv.Atoi(offsetPar)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	categoryPar, exist := ctx.GetQuery("categoryId")
	if exist && categoryPar != "" {
		categoryId, err = strconv.Atoi(categoryPar)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	response := control.svc.FindAllCampaign(ctx, ctx.Query("order"), categoryId, ctx.Query("search"), offset)
	WriteJsonResponse(ctx, response)
}

func (control *CampaignController) FindCampaignByUser(ctx *gin.Context) {
	response := control.svc.FindCampaignByUser(ctx, ctx.Param("userAddress"))
	WriteJsonResponse(ctx, response)
}

func (control *CampaignController) FindCampaignByAddress(ctx *gin.Context) {
	response := control.svc.FindCampaignByAddress(ctx, ctx.Param("address"))
	WriteJsonResponse(ctx, response)
}

func (control *CampaignController) CreateCampaign(ctx *gin.Context) {
	var req params.CreateCampaign
	err := ctx.ShouldBind(&req)
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

	response := control.svc.CreateCampaign(ctx, &req)
	WriteJsonResponse(ctx, response)
}

func (control *CampaignController) CreateProposal(ctx *gin.Context) {
	var req params.CreateProposal
	err := ctx.ShouldBind(&req)
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

	response := control.svc.CreateProposal(ctx, &req)
	WriteJsonResponse(ctx, response)
}

func (control *CampaignController) FindProposalByAddress(ctx *gin.Context) {
	response := control.svc.FindProposalByAddress(ctx, ctx.Param("address"))
	WriteJsonResponse(ctx, response)
}
