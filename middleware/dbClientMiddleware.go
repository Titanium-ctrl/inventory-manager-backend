package middleware

import (
	"github.com/gofiber/fiber/v2"
	"ucrs.com/inventory-manager/backend/internal/database"
)

func DBClientMiddleware(c *fiber.Ctx) error {
	//Fetch JWT from request header
	jwt := c.Get("Authorization")
	if jwt == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	supabaseClient := database.CreateClient(jwt)
	c.Locals("supabaseClient", supabaseClient)

	return c.Next()
}
