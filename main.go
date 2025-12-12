package main

import (
	"log"

	"prestasi_mhs/app/repository"
	"prestasi_mhs/app/service"
	"prestasi_mhs/config"
	"prestasi_mhs/database"
	"prestasi_mhs/route"

	fiberSwagger "github.com/swaggo/fiber-swagger"
	_"prestasi_mhs/docs"
)

// @title Prestasi Mahasiswa API
// @version 1.0
// @description REST API untuk Sistem Prestasi Mahasiswa
// @contact.name Shendy Tria Amelyana
// @contact.email shendytriaaa@gmail.com
// @host localhost:3000
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
	lecturerSvc := service.NewLecturerService(lecturerRepo, studentRepo)

	authSvc := service.NewAuthService(userRepo)
	userSvc := service.NewUserService(userRepo)
	achSvc := service.NewAchievementService(achRepo, studentSvc)
	reportSvc := service.NewReportService(reportRepo, studentRepo, lecturerRepo)

	app := config.NewApp()
	
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	route.RegisterRoutes(app, authSvc, userSvc, achSvc, studentSvc, lecturerSvc, reportSvc)

	log.Println("Server running on port", config.C.AppPort)
	log.Fatal(app.Listen(":" + config.C.AppPort))
}
