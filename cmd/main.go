package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nathanramli/solcare-backend/config"
	"github.com/nathanramli/solcare-backend/httpserver"
	"github.com/nathanramli/solcare-backend/httpserver/controllers"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/gorm"
	"github.com/nathanramli/solcare-backend/httpserver/services"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("cannot load your env")
	}
}

func main() {
	db, err := config.ConnectPostgresGORM()
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	config.GenerateJwtSignature()

	//repo
	userRepo := gorm.NewUserRepo(db)
	campaignRepo := gorm.NewCampaignRepo(db)
	categoryRepo := gorm.NewCategoryRepo(db)
	proposalRepo := gorm.NewProposalRepo(db)
	reportRepo := gorm.NewReportRepo(db)
	kycQueueRepo := gorm.NewKyqQueueRepo(db)
	adminRepo := gorm.NewAdminRepo(db)

	userSvc := services.NewUserSvc(userRepo, kycQueueRepo, adminRepo)
	userHandler := controllers.NewUserController(userSvc)

	campaignSvc := services.NewCampaignSvc(campaignRepo, categoryRepo, proposalRepo)
	campaignHandler := controllers.NewCampaignController(campaignSvc)

	categorySvc := services.NewCategorySvc(categoryRepo)
	categoryHandler := controllers.NewCategoryController(categorySvc)

	reportSvc := services.NewReportSvc(reportRepo, campaignRepo)
	reportHandler := controllers.NewReportController(reportSvc)

	app := httpserver.NewRouter(router, userHandler, campaignHandler, categoryHandler, reportHandler)
	PORT := os.Getenv("PORT")
	app.Start(":" + PORT)
}
