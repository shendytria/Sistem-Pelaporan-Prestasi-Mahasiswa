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
	api.Post("/auth/login", authSvc.Login)
	api.Post("/auth/refresh", middleware.JWT(), authSvc.Refresh)
	api.Get("/auth/profile", middleware.JWT(), authSvc.Profile)
	api.Post("/auth/logout", middleware.JWT(), authSvc.Logout)

	users := api.Group("/users", middleware.JWT(), middleware.Permission("manage_users"))
	users.Get("/", userSvc.List)
	users.Get("/:id", userSvc.Detail)
	users.Post("/", userSvc.Create)
	users.Put("/:id", userSvc.Update)
	users.Delete("/:id", userSvc.Delete)
	users.Put("/:id/role", userSvc.UpdateRole)

	ach := api.Group("/achievements", middleware.JWT())
	ach.Get("/", middleware.Permission("read_achievement"), achSvc.List)
	ach.Get("/:id", middleware.Permission("read_achievement"), achSvc.Detail)
	ach.Post("/", middleware.Permission("create_achievement"), achSvc.Create)
	ach.Put("/:id", middleware.Permission("update_achievement"), achSvc.Update)
	ach.Delete("/:id", middleware.Permission("delete_achievement"), achSvc.Delete)
	ach.Post("/:id/submit", middleware.Permission("update_achievement"), achSvc.Submit)
	ach.Post("/:id/verify", middleware.Permission("verify_achievement"), achSvc.Verify)
	ach.Post("/:id/reject", middleware.Permission("verify_achievement"), achSvc.Reject)
	ach.Get("/:id/history", middleware.Permission("read_achievement"), achSvc.History)
	ach.Post("/:id/attachments", middleware.Permission("update_achievement"), achSvc.AddAttachment)

	students := api.Group("/students", middleware.JWT())
	students.Get("/", middleware.Permission("read_achievement"), studentSvc.List)
	students.Get("/:id", middleware.Permission("read_achievement"), studentSvc.Detail)
	students.Get("/:id/achievements", middleware.Permission("read_achievement"), studentSvc.ListByStudent)
	students.Put("/:id/advisor", middleware.Permission("manage_users"), studentSvc.UpdateAdvisor)

	lect := api.Group("/lecturers", middleware.JWT())
	lect.Get("/", middleware.Permission("read_achievement"), lecturerSvc.List)
	lect.Get("/:id/advisees", middleware.Permission("verify_achievement"), lecturerSvc.Advisees)

	reports := api.Group("/reports", middleware.JWT())
	reports.Get("/statistics", middleware.Permission("read_achievement"), reportSvc.Statistics)
	reports.Get("/student/:id", middleware.Permission("read_achievement"), reportSvc.StudentReport)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Prestasi Mahasiswa API is running",
		})
	})
}
