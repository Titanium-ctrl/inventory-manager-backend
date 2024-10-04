package routes

import (
	"ucrs.com/inventory-manager/backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Product routes
	app.Post("/products", handlers.CreateProduct)
	app.Get("/products", handlers.GetProducts)
	app.Get("/products/:id", handlers.GetProduct)
	app.Put("/products/:id", handlers.UpdateProduct)
	app.Delete("/products/:id", handlers.DeleteProduct)

	// Attribute routes
	app.Post("/attributes", handlers.CreateAttribute)
	app.Put("/attributes/:id", handlers.UpdateAttribute)
	app.Get("/attributes", handlers.GetAttributes)
	app.Get("/attributes/:id", handlers.GetAttribute)
	app.Delete("/attributes/:id", handlers.DeleteAttribute)

	// Other routes...
}
