package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
	"ucrs.com/inventory-manager/backend/internal/database"
	"ucrs.com/inventory-manager/backend/internal/models"
)

func CreateProduct(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	product := new(models.Product)

	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Basic validation
	if product.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product name is required",
		})
	}

	if product.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Price must be greater than 0",
		})
	}

	// Set timestamps
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	// Set UserID and ID
	product.UserID = database.FetchUserID(supabaseClient)
	product.ID = uuid.New()

	//Save to database
	result, count, err := supabaseClient.From("products").Insert(product, false, "", "", "").Execute() //I believe the other params are correct
	fmt.Println(string(result), count, err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save product to database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// We'll implement other handlers (GetProducts, GetProduct, UpdateProduct, DeleteProduct) later
