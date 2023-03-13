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
	report   *controllers.ReportController
}

func NewRouter(r *gin.Engine, user *controllers.UserController, campaign *controllers.CampaignController, category *controllers.CategoryController, report *controllers.ReportController) *router {
	return &router{
		router:   r,
		user:     user,
		campaign: campaign,
		category: category,
		report:   report,
	}
}

func (r *router) Start(port string) {
	r.router.Use(cors)

	r.router.Static("/resources/", "./resources")

	r.router.Use(contentJson)

	r.router.POST("/v1/users/login", r.user.Login)
	r.router.PUT("/v1/users/info/:address", r.verifyToken, r.user.UpdateUser)
	r.router.PUT("/v1/users/avatar/:address", r.verifyToken, r.user.UpdateAvatar)
	r.router.GET("/v1/users/info/:address", r.user.FindUserByAddress)
	r.router.GET("/v1/users", r.user.FindAllUsers)
	r.router.POST("/v1/users/kyc", r.verifyToken, r.user.RequestKyc)
	r.router.GET("/v1/users/kyc", r.verifyToken, r.user.FindKycRequestByUser)

	r.router.POST("/v1/report", r.verifyToken, r.report.CreateReport)
	r.router.GET("/v1/report/:id", r.report.FindReportById)

	r.router.POST("/v1/campaign", r.campaign.CreateCampaign)
	r.router.GET("/v1/campaign/user/:userAddress", r.campaign.FindCampaignByUser)
	r.router.POST("/v1/campaign/proposal", r.campaign.CreateProposal)
	r.router.GET("/v1/campaign/proposal/:address", r.campaign.FindProposalByAddress)
	r.router.GET("/v1/campaign", r.campaign.FindAllCampaign)
	r.router.GET("/v1/campaign/:address", r.campaign.FindCampaignByAddress)
	r.router.POST("/v1/campaign/evidence", r.verifyToken, r.campaign.UploadEvidence)
	r.router.GET("/v1/campaign/evidence", r.verifyAdminToken, r.campaign.FindAllCampaignWithEvidence)
	r.router.POST("/v1/campaign/evidence/verify", r.verifyAdminToken, r.campaign.VerifyEvidence)

	r.router.GET("/v1/categories", r.category.FindAllCategories)
	r.router.GET("/v1/categories/:categoryId", r.category.FindCategoryById)

	r.router.GET("/v1/admins/kyc", r.verifyAdminToken, r.user.FindAllKycRequest)
	r.router.POST("/v1/admins/kyc", r.verifyAdminToken, r.user.VerifyKyc)
	r.router.DELETE("/v1/admins/kyc/:address", r.verifyAdminToken, r.user.RemoveKyc)

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

func (r *router) verifyAdminToken(ctx *gin.Context) {
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

	if !claims.IsAdmin {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "you are not an admin",
		})
		return
	}

	ctx.Set("userData", claims)
}

func cors(c *gin.Context) {
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

func contentJson(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	c.Next()
}
