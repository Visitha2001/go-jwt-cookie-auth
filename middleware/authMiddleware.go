package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/visitha2001/go-jwt-auth/configs"
)

// AuthRequired is a middleware to protect routes
func AuthRequired(c *fiber.Ctx) error {
	// Get token from cookie
	token := c.Cookies("jwt")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing or malformed JWT",
		})
	}

	// Validate token
	claims, err := configs.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid or expired JWT",
		})
	}

	// Set user ID in locals for the next handler
	c.Locals("userID", claims.UserID)

	return c.Next()
}
