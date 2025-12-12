package middleware

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"prestasi_mhs/utils"
	"strings"

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

		if v, ok := claims["role_id"].(string); ok {
			c.Locals("role_id", v)
		}

		if v, ok := claims["permissions"].([]interface{}); ok {
			perms := make([]string, 0)
			for _, p := range v {
				if s, ok := p.(string); ok {
					perms = append(perms, s)
				}
			}
			c.Locals("permissions", perms)
		}

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

func HasPermission(c *fiber.Ctx, permission string) bool {
	perms, ok := c.Locals("permissions").([]string)
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == "*" {
			return true
		}
	}
	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}

func Permission(requiredPermission string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		perms, ok := c.Locals("permissions").([]string)
		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"error": "no permissions in token",
			})
		}

		for _, p := range perms {
			if p == "*" {
				return c.Next()
			}
		}

		for _, p := range perms {
			if p == requiredPermission {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "forbidden – missing permission: " + requiredPermission,
		})
	}
}
