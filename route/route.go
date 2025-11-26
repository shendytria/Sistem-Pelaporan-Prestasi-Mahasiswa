package route

import (
	"prestasi_mhs/app/service"
	"prestasi_mhs/constant"
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
) {

	api := app.Group("/api")

	api.Post("/auth/login", authSvc.LoginHTTP)

	users := api.Group("/users", middleware.JWT(), middleware.Role(constant.RoleAdmin))
	users.Get("/", userSvc.ListHTTP)
	users.Get("/:id", userSvc.DetailHTTP)
	users.Post("/", userSvc.CreateHTTP)
	users.Put("/:id", userSvc.UpdateHTTP)
	users.Delete("/:id", userSvc.DeleteHTTP)
	users.Put("/:id/role", userSvc.UpdateRoleHTTP)

	ach := api.Group("/achievements", middleware.JWT())
	ach.Get("/me", middleware.Role(constant.RoleMahasiswa), achSvc.ListMineHTTP)
	ach.Get("/", middleware.Role(constant.RoleAdmin), achSvc.ListAllHTTP)
	ach.Get("/:id", middleware.Role(constant.RoleMahasiswa, constant.RoleDosenWali, constant.RoleAdmin), achSvc.DetailHTTP)
	ach.Post("/", middleware.Role(constant.RoleMahasiswa), achSvc.CreateHTTP)
	ach.Put("/:id", middleware.Role(constant.RoleMahasiswa), achSvc.UpdateHTTP)
	ach.Delete("/:id", middleware.Role(constant.RoleMahasiswa, constant.RoleAdmin), achSvc.DeleteHTTP)
	ach.Post("/:id/submit", middleware.Role(constant.RoleMahasiswa, constant.RoleAdmin), achSvc.SubmitHTTP)
	ach.Post("/:id/verify", middleware.Role(constant.RoleDosenWali, constant.RoleAdmin), achSvc.VerifyHTTP)
	ach.Post("/:id/reject", middleware.Role(constant.RoleDosenWali, constant.RoleAdmin), achSvc.RejectHTTP)
	ach.Get("/:id/history", middleware.Role(constant.RoleMahasiswa, constant.RoleDosenWali, constant.RoleAdmin), achSvc.HistoryHTTP)
	ach.Post("/:id/attachments", middleware.Role(constant.RoleMahasiswa, constant.RoleAdmin), achSvc.AddAttachmentHTTP)

	students := api.Group("/students", middleware.JWT())
	students.Get("/", middleware.Role(constant.RoleAdmin), studentSvc.ListHTTP)
	students.Get("/me", middleware.Role(constant.RoleMahasiswa), studentSvc.MeHTTP)
	students.Get("/:id", middleware.Role(constant.RoleAdmin, constant.RoleDosenWali), studentSvc.DetailHTTP)
	students.Get("/:id/achievements", middleware.Role(constant.RoleAdmin, constant.RoleDosenWali), achSvc.ListByStudentHTTP)
	students.Put("/:id/advisor", middleware.Role(constant.RoleAdmin), studentSvc.UpdateAdvisorHTTP)

	lect := api.Group("/lecturers", middleware.JWT())
	lect.Get("/", middleware.Role(constant.RoleAdmin), lecturerSvc.ListHTTP)
	lect.Get("/:id/advisees", middleware.Role(constant.RoleAdmin, constant.RoleDosenWali), lecturerSvc.AdviseesHTTP)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Prestasi Mahasiswa API is running",
		})
	})
}
