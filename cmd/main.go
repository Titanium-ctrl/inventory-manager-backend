package main

import (
	"github.com/gofiber/fiber/v2"
	"ucrs.com/inventory-manager/backend/internal/database"
	"ucrs.com/inventory-manager/backend/internal/routes"
)

func main() {
	app := fiber.New()

	database.Connect()
	routes.SetupRoutes(app)

	app.Listen(":3000")
}
