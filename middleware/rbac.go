package middleware

import "github.com/gofiber/fiber/v2"

func Role(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		roleID, ok := c.Locals("role_id").(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token data"})
		}

		for _, r := range allowedRoles {
			if roleID == r {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
	}
}
