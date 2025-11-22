package route

import (
	"prestasi_mhs/app/service"
	"prestasi_mhs/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	authSvc *service.AuthService,
	achSvc *service.AchievementUsecaseService,
) {

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/login", authSvc.LoginHTTP)

	ach := api.Group("/achievements", middleware.JWT())

	ach.Post("/", achSvc.CreateHTTP)
	ach.Get("/:studentId", achSvc.ListByStudentHTTP)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Prestasi Mahasiswa API is running 🚀",
		})
	})
}
