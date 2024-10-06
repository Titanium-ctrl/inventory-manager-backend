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

	// SKU routes
	app.Post("/skus", handlers.CreateSKU)
	app.Put("/skus/:id", handlers.UpdateSKU)
	app.Get("/skus/:id", handlers.GetSKU)
	app.Get("/skus", handlers.GetSKUs)
	app.Get("/skus/:id/products", handlers.GetSKUsByProductID)
	app.Delete("/skus/:id", handlers.DeleteSKU)

	// Barcode routes
	app.Post("/barcodes", handlers.CreateBarcode)
	app.Put("/barcodes/:id", handlers.UpdateBarcode)
	app.Get("/barcodes", handlers.GetBarcodes)
	app.Get("/barcodes/:id", handlers.GetBarcode)
	app.Get("/barcodes/:id/skus", handlers.GetBarcodesBySKUID)
	app.Delete("/barcodes/:id", handlers.DeleteBarcode)

	//Category routes
	app.Post("/categories", handlers.CreateCategory)
	app.Put("/categories/:id", handlers.UpdateCategory)
	app.Get("/categories", handlers.GetCategories)
	app.Get("/categories/:id", handlers.GetCategory)
	app.Delete("/categories/:id", handlers.DeleteCategory)
	app.Get("/categories/:id/parent", handlers.GetCategoriesByParentID)

	// Other routes...
}
