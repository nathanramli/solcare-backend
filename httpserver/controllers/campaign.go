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
	var offset = 0

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

	response := control.svc.FindAllCampaign(ctx, ctx.Query("order"), offset)
	WriteJsonResponse(ctx, response)
}

func (control *CampaignController) FindCampaignByUser(ctx *gin.Context) {
	response := control.svc.FindCampaignByUser(ctx, ctx.Param("userAddress"))
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
