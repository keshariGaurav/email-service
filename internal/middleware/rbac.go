package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequireRole creates a middleware that checks if the user has the required role
func RequireRole(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*JWTClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid authentication",
			})
		}

		// Check if user's role is in the allowed roles
		for _, role := range allowedRoles {
			if claims.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}