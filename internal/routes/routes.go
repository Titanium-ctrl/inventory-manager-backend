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

	// SKU Attribute routes
	app.Post("/sku/:skuid/attributes", handlers.UpdateSKUAttribute)       //Insert/update skuattribute
	app.Get("/sku/:skuid/attributes", handlers.GetSKUAttributes)          //Get attributes for a sku
	app.Get("/sku/:skuid/attributes/:id", handlers.GetSKUAttribute)       //get specific sku attributw
	app.Delete("/sku/:skuid/attributes/:id", handlers.DeleteSKUAttribute) //delete a specific attribute for a sku

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

	// Location (warehouses) routes
	app.Post("/warehouses", handlers.CreateWarehouse)
	app.Put("/warehouses/:id", handlers.UpdateWarehouse)
	app.Get("/warehouses", handlers.GetWarehouses)
	app.Get("/warehouses/:id", handlers.GetWarehouse)
	app.Delete("/warehouses/:id", handlers.DeleteWarehouse)

	//Inventory routes - CRUD functions for database table storing quantity of items in inventory
	app.Post("/inventory/:locationid/:skuid", handlers.UpdateInventory)         //Add/update inventory quantity
	app.Get("/inventory", handlers.GetInventory)                                //List locations only
	app.Get("/inventory/:locationid", handlers.GetInventory)                    //Get products stored at said location
	app.Get("/inventory/:locationid/sku/:skuid", handlers.GetSpecificInventory) //Get quantity of specific sku at specific location
	app.Get("/inventory/sku/:skuid", handlers.GetInventoryForSKU)               //Get quantity of specific sku at all locations
	app.Delete("/inventory/:locationid", handlers.DeleteInventory)              //Delete inventory location (may not be needed)

	//User details routes
	app.Post("/users", handlers.CreateUser)
	app.Get("/users/company/:companyid", handlers.GetUsersFromCompanyID)
	app.Get("/users", handlers.GetUser)
	app.Put("/users/:id", handlers.UpdateUser)
	app.Delete("/users/:id", handlers.DeleteUser) // Consider only admin users to be able to delete users

	//Company routes
	app.Post("/companies", handlers.CreateCompany)
	//app.Get("/companies", handlers.GetCompanies)
	app.Get("/companies/:id", handlers.GetCompany)
	app.Put("/companies/:id", handlers.UpdateCompany)
	app.Delete("/companies/:id", handlers.DeleteCompany)
}
