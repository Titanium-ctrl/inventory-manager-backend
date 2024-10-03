package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
	"ucrs.com/inventory-manager/backend/internal/database"
	"ucrs.com/inventory-manager/backend/internal/models"
	"ucrs.com/inventory-manager/backend/pkg"
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
	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Cannot fetch user ID - please log in",
		})
	}
	product.UserID = userID
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

// Function to fetch multiple products, split into pages of max 10 products
func GetProducts(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	startIndex, endIndex := pkg.GetPaginationIndexes(page, 10)
	products, _, err := supabaseClient.From("products").Select("*", "", false).Range(startIndex, endIndex, "").Execute()
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

	return c.Status(fiber.StatusOK).JSON(respStruct)
}

// Function to fetch a single product, based on it's ID which should be passsed as a param
func GetProduct(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	productID := c.Params("id")
	product, _, err := supabaseClient.From("products").Select("*", "", false).Eq("id", productID).Execute()
	//fmt.Println(string(product), err)
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

// Function to update a product, based on its ID which should be passsed as a param
func UpdateProduct(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	productID := c.Params("id")
	product := new(models.Product)
	pid, err := uuid.Parse(productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}
	product.ID = pid
	product.UserID = database.FetchUserID(supabaseClient)

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
	_, _, err = supabaseClient.From("products").Update(product, "", "").Eq("id", productID).Execute()
	//fmt.Println(string(result), count, err)
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
	_, _, err := supabaseClient.From("products").Delete("", "").Eq("id", productID).Execute()
	//fmt.Println(string(result), count, err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete product from database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}
