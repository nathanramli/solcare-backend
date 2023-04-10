package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nathanramli/solcare-backend/common"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/params"
	"github.com/nathanramli/solcare-backend/httpserver/services"
	"net/http"
)

type TransactionController struct {
	svc services.TransactionSvc
}

func NewTransactionController(svc services.TransactionSvc) *TransactionController {
	return &TransactionController{
		svc: svc,
	}
}

func (control *TransactionController) CreateTransaction(ctx *gin.Context) {
	var req params.CreateTransaction
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

	response := control.svc.CreateTransaction(ctx, userData.Address, &req)
	WriteJsonResponse(ctx, response)
}

func (control *TransactionController) FindAllTransactionsByUser(ctx *gin.Context) {
	response := control.svc.FindAllTransactionsByUser(ctx, ctx.Param("address"))
	WriteJsonResponse(ctx, response)
}
