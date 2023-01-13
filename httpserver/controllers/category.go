package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nathanramli/solcare-backend/httpserver/services"
	"net/http"
	"strconv"
)

type CategoryController struct {
	svc services.CategoriesSvc
}

func NewCategoryController(svc services.CategoriesSvc) *CategoryController {
	return &CategoryController{
		svc: svc,
	}
}

func (control *CategoryController) FindCategoryById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("categoryId"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := control.svc.FindCategoryById(ctx, uint(id))
	WriteJsonResponse(ctx, response)
}

func (control *CategoryController) FindAllCategories(ctx *gin.Context) {
	response := control.svc.FindAllCategories(ctx)
	WriteJsonResponse(ctx, response)
}
