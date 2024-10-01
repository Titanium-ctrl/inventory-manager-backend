package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"ucrs.com/inventory-manager/backend/internal/routes"
	"ucrs.com/inventory-manager/backend/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Use(middleware.DBClientMiddleware)

	routes.SetupRoutes(app)

	app.Listen(":3000")
}
