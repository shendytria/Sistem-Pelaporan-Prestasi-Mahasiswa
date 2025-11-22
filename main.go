package main

import (
	"log"

	"prestasi_mhs/app/repository"
	"prestasi_mhs/app/service"
	"prestasi_mhs/config"
	"prestasi_mhs/database"
	"prestasi_mhs/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	config.Load()

	database.ConnectPostgres()
	database.ConnectMongo()

	userRepo := repository.NewUserRepository()
	studentRepo := repository.NewStudentRepository()
	achMongoRepo := repository.NewAchievementMongoRepository()
	achRefRepo := repository.NewAchievementReferenceRepository()

	authSvc := service.NewAuthService(userRepo)
	achUsecase := service.NewAchievementUsecaseService(
		achMongoRepo,
		achRefRepo,
		studentRepo,
	)

	app := fiber.New()

	app.Use(logger.New())

	route.RegisterRoutes(app, authSvc, achUsecase)

	log.Println("Server running on port", config.C.AppPort)
	log.Fatal(app.Listen(":" + config.C.AppPort))
}
