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
	achRepo := repository.NewAchievementRepository()
	lecturerRepo := repository.NewLecturerRepository()
	reportRepo := repository.NewReportRepository()

	studentSvc := service.NewStudentService(studentRepo, achRepo)
	lecturerSvc := service.NewLecturerService(lecturerRepo)

	authSvc := service.NewAuthService(userRepo)
	userSvc := service.NewUserService(userRepo)
	achSvc := service.NewAchievementService(achRepo, studentSvc)
	reportSvc := service.NewReportService(reportRepo, studentRepo, lecturerRepo)

	app := fiber.New()

	app.Use(logger.New())

	route.RegisterRoutes(app, authSvc, userSvc, achSvc, studentSvc, lecturerSvc, reportSvc)

	log.Println("Server running on port", config.C.AppPort)
	log.Fatal(app.Listen(":" + config.C.AppPort))
}
