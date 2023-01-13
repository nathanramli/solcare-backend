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

	userSvc := services.NewUserSvc(userRepo)
	userHandler := controllers.NewUserController(userSvc)

	campaignSvc := services.NewCampaignSvc(campaignRepo, categoryRepo)
	campaignHandler := controllers.NewCampaignController(campaignSvc)

	categorySvc := services.NewCategorySvc(categoryRepo)
	categoryHandler := controllers.NewCategoryController(categorySvc)

	app := httpserver.NewRouter(router, userHandler, campaignHandler, categoryHandler)
	PORT := os.Getenv("PORT")
	app.Start(":" + PORT)
}
