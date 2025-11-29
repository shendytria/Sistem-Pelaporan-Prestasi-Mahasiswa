package route

import (
	"prestasi_mhs/app/service"
	"prestasi_mhs/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	authSvc *service.AuthService,
	userSvc *service.UserService,
	achSvc *service.AchievementUsecaseService,
	studentSvc *service.StudentService,
	lecturerSvc *service.LecturerService,
	reportSvc *service.ReportService,
) {

	api := app.Group("/api")
	api.Post("/auth/login", authSvc.LoginHTTP)
	api.Post("/auth/refresh", middleware.JWT(), authSvc.RefreshHTTP)
	api.Get("/auth/profile", middleware.JWT(), authSvc.ProfileHTTP)
	api.Post("/auth/logout", middleware.JWT(), authSvc.LogoutHTTP)

	users := api.Group("/users", middleware.JWT(), middleware.Role("Admin"))
	users.Get("/", userSvc.ListHTTP)
	users.Get("/:id", userSvc.DetailHTTP)
	users.Post("/", userSvc.CreateHTTP)
	users.Put("/:id", userSvc.UpdateHTTP)
	users.Delete("/:id", userSvc.DeleteHTTP)
	users.Put("/:id/role", userSvc.UpdateRoleHTTP)

	ach := api.Group("/achievements", middleware.JWT())
	ach.Get("/", middleware.Role("Admin", "Mahasiswa", "Dosen Wali"), achSvc.ListHTTP)
	ach.Get("/:id", middleware.Role("Mahasiswa", "Dosen Wali", "Admin"), achSvc.DetailHTTP)
	ach.Post("/", middleware.Role("Mahasiswa", "Admin"), achSvc.CreateHTTP)
	ach.Put("/:id", middleware.Role("Mahasiswa", "Admin"), achSvc.UpdateHTTP)
	ach.Delete("/:id", middleware.Role("Mahasiswa", "Admin"), achSvc.DeleteHTTP)
	ach.Post("/:id/submit", middleware.Role("Mahasiswa", "Admin"), achSvc.SubmitHTTP)
	ach.Post("/:id/verify", middleware.Role("Dosen Wali", "Admin"), achSvc.VerifyHTTP)
	ach.Post("/:id/reject", middleware.Role("Dosen Wali", "Admin"), achSvc.RejectHTTP)
	ach.Get("/:id/history", middleware.Role("Mahasiswa", "Dosen Wali", "Admin"), achSvc.HistoryHTTP)
	ach.Post("/:id/attachments", middleware.Role("Mahasiswa", "Admin"), achSvc.AddAttachmentHTTP)

	students := api.Group("/students", middleware.JWT())
	students.Get("/", middleware.Role("Admin", "Dosen Wali"), studentSvc.ListHTTP)
	students.Get("/:id", middleware.Role("Admin", "Dosen Wali"), studentSvc.DetailHTTP)
	students.Get("/:id/achievements", middleware.Role("Admin", "Dosen Wali"), studentSvc.ListByStudentHTTP)
	students.Put("/:id/advisor", middleware.Role("Admin"), studentSvc.UpdateAdvisorHTTP)

	lect := api.Group("/lecturers", middleware.JWT())
	lect.Get("/", middleware.Role("Admin", "Mahasiswa"), lecturerSvc.ListHTTP)
	lect.Get("/:id/advisees", middleware.Role("Admin", "Dosen Wali"), lecturerSvc.AdviseesHTTP)

	reports := api.Group("/reports", middleware.JWT())
	reports.Get("/statistics", middleware.Role("Admin", "Dosen Wali", "Mahasiswa"), reportSvc.StatisticsHTTP,)
	reports.Get("/student/:id", middleware.Role("Admin", "Dosen Wali", "Mahasiswa"), reportSvc.StudentReportHTTP,)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Prestasi Mahasiswa API is running",
		})
	})
}
