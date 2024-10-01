package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
	"ucrs.com/inventory-manager/backend/internal/database"
	"ucrs.com/inventory-manager/backend/internal/models"
)

// Create a new product, based on the JSON passed in the body
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

// ADD PAGING TO COMPLETE THIS FUNCTION, and test it
func GetProducts(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	products, _, err := supabaseClient.From("products").Select("*", "", false).Execute()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch products from database",
		})
	}

	respStruct := []struct {
		models.Product
	}{}

	err = json.Unmarshal(products, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal product from database",
		})

	}

	if len(respStruct) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Products not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

// Function to fetch a single product, based on it's ID which should be passsed as a param
func GetProduct(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	productID := c.Params("id")
	product, _, err := supabaseClient.From("products").Select("*", "", false).Eq("id", productID).Execute()
	fmt.Println(string(product), err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch product from database",
		})
	}

	respStruct := []struct {
		models.Product
	}{}

	err = json.Unmarshal(product, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal product from database",
		})

	}

	if len(respStruct) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(respStruct[0])
}

// Work on finishing and testing this function
func UpdateProduct(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	productID := c.Params("id")
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
	product.UpdatedAt = now

	//Save to database
	result, count, err := supabaseClient.From("products").Update(product, "", "").Eq("id", productID).Execute()
	fmt.Println(string(result), count, err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save product to database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

// Work on finishing and testing this function
func DeleteProduct(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	productID := c.Params("id")

	//Save to database
	result, count, err := supabaseClient.From("products").Delete("", "").Eq("id", productID).Execute()
	fmt.Println(string(result), count, err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete product from database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}
