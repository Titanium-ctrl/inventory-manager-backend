package routes

import (
	"ucrs.com/inventory-manager/backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Product routes
	app.Post("/products", handlers.CreateProduct)
	// Other routes...
}
