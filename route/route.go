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
	achSvc *service.AchievementUsecaseService,
) {

	api := app.Group("/api")

	api.Post("/auth/login", authSvc.LoginHTTP)

	ach := api.Group("/achievements", middleware.JWT())

	ach.Post("/",
		middleware.Role(constant.RoleMahasiswa),
		achSvc.CreateHTTP,
	)
	ach.Get("/me",
		middleware.Role(constant.RoleMahasiswa),
		achSvc.ListMineHTTP,
	)

	ach.Get("/student/:studentId",
		middleware.Role(constant.RoleDosenWali, constant.RoleAdmin),
		achSvc.ListByStudentHTTP,
	)

	ach.Get("/",
		middleware.Role(constant.RoleAdmin),
		achSvc.ListAllHTTP,
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Prestasi Mahasiswa API is running 🚀",
		})
	})
}
