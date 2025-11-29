package middleware

import (
	"prestasi_mhs/utils"
	"strings"
	"errors"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
)

func JWT() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		bearer := c.Get("Authorization")
		if bearer == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing token"})
		}

		tokenStr := strings.TrimPrefix(bearer, "Bearer ")
		token, claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return c.Status(401).JSON(fiber.Map{
					"error": "token expired",
				})
			}

			return c.Status(401).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		if !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}

		if v, ok := claims["user_id"].(string); ok {
			c.Locals("user_id", v)
		} else {
			c.Locals("user_id", "")
		}

		if v, ok := claims["role"].(string); ok {
			c.Locals("role", v)
		}

		c.Locals("permissions", claims["permissions"])

		return c.Next()
	}
}

func Role(allowedRoles ...string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid token data",
			})
		}

		for _, r := range allowedRoles {
			if role == r {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "forbidden",
		})
	}
}
