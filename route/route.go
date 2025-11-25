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
	ach.Post("/", middleware.Role(constant.RoleMahasiswa), achSvc.CreateHTTP)
	ach.Get("/me", middleware.Role(constant.RoleMahasiswa), achSvc.ListMineHTTP)
	ach.Get("/", middleware.Role(constant.RoleAdmin), achSvc.ListAllHTTP)
	ach.Get("/:id", middleware.Role(constant.RoleMahasiswa, constant.RoleDosenWali, constant.RoleAdmin), achSvc.DetailHTTP)
	ach.Get("/:id/history", middleware.Role(constant.RoleMahasiswa, constant.RoleDosenWali, constant.RoleAdmin), achSvc.HistoryHTTP)
	ach.Post("/:id/attachments", middleware.Role(constant.RoleMahasiswa, constant.RoleAdmin), achSvc.AddAttachmentHTTP)
	ach.Put("/:id", middleware.Role(constant.RoleMahasiswa), achSvc.UpdateHTTP)
	ach.Delete("/:id", middleware.Role(constant.RoleMahasiswa, constant.RoleAdmin), achSvc.DeleteHTTP)
	ach.Post("/:id/submit", middleware.Role(constant.RoleMahasiswa, constant.RoleAdmin), achSvc.SubmitHTTP)
	ach.Post("/:id/verify", middleware.Role(constant.RoleDosenWali, constant.RoleAdmin), achSvc.VerifyHTTP)
	ach.Post("/:id/reject", middleware.Role(constant.RoleDosenWali, constant.RoleAdmin), achSvc.RejectHTTP)

	// harusnya bukan achsvc
	// students := api.Group("/students", middleware.JWT())
	// students.Get("/:id/achievements", middleware.Role(constant.RoleDosenWali, constant.RoleAdmin), achSvc.ListByStudentHTTP)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Prestasi Mahasiswa API is running",
		})
	})
}
