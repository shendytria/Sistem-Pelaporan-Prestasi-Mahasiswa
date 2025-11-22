package middleware

import (
	"prestasi_mhs/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {

		bearer := c.Get("Authorization")
		if bearer == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing token"})
		}

		tokenStr := strings.TrimPrefix(bearer, "Bearer ")

		token, claims, err := utils.ParseJWT(tokenStr)
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("user_id", claims["user_id"])
		c.Locals("role_id", claims["role_id"])

		return c.Next()
	}
}
