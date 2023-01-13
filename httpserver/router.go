package httpserver

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nathanramli/solcare-backend/common"
	"github.com/nathanramli/solcare-backend/httpserver/controllers"
)

type router struct {
	router *gin.Engine

	user     *controllers.UserController
	campaign *controllers.CampaignController
	category *controllers.CategoryController
}

func NewRouter(r *gin.Engine, user *controllers.UserController, campaign *controllers.CampaignController, category *controllers.CategoryController) *router {
	return &router{
		router:   r,
		user:     user,
		campaign: campaign,
		category: category,
	}
}

func (r *router) Start(port string) {
	r.router.Use(cors)
	// users
	r.router.POST("/v1/users/login", r.user.Login)

	r.router.POST("/v1/campaign", r.campaign.CreateCampaign)

	r.router.GET("/v1/categories", r.category.FindAllCategories)
	r.router.GET("/v1/categories/:categoryId", r.category.FindCategoryById)

	r.router.Run(port)
}

func (r *router) verifyToken(ctx *gin.Context) {
	bearerToken := strings.Split(ctx.Request.Header.Get("Authorization"), "Bearer ")
	if len(bearerToken) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid bearer token",
		})
		return
	}
	claims, err := common.ValidateToken(bearerToken[1])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Set("userData", claims)
}

func cors(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
