package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/params"
	"github.com/nathanramli/solcare-backend/httpserver/services"
)

type UserController struct {
	svc services.UserSvc
}

func NewUserController(svc services.UserSvc) *UserController {
	return &UserController{
		svc: svc,
	}
}

func (control *UserController) Login(ctx *gin.Context) {
	var req params.Login
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

	response := control.svc.Login(ctx, &req)
	WriteJsonResponse(ctx, response)
}

func (control *UserController) UpdateUser(ctx *gin.Context) {
	var req params.UpdateUser
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

	response := control.svc.UpdateUser(ctx, ctx.Param("address"), &req)
	WriteJsonResponse(ctx, response)
}

func (control *UserController) FindUserByAddress(ctx *gin.Context) {
	response := control.svc.FindUserByAddress(ctx, ctx.Param("address"))
	WriteJsonResponse(ctx, response)
}
