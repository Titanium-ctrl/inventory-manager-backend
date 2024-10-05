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

func GetSKU(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)

	skuID := c.Params("id")
	sku, _, err := supabaseClient.From("skus").Select("*", "", false).Eq("id", skuID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch SKU from database",
		})
	}

	respStruct := []struct {
		models.SKU
	}{}

	err = json.Unmarshal(sku, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal SKU from database",
		})

	}

	if len(respStruct) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "SKU not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(respStruct[0])
}

// DO I LEAVE IT LIKE THIS OR DO I SET IT TO RETURN ALL SKUS BASED ON THE PRODUCT ID???
func GetSKUs(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	skus, _, err := supabaseClient.From("skus").Select("*", "", false).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch SKUs from database",
		})
	}
	respStruct := []struct {
		models.SKU
	}{}
	err = json.Unmarshal(skus, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal SKU from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct)
}

func GetSKUsByProductID(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	productID := c.Params("id")
	skus, _, err := supabaseClient.From("skus").Select("*", "", false).Eq("product_id", productID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch SKUs from database",
		})
	}
	respStruct := []struct {
		models.SKU
	}{}
	err = json.Unmarshal(skus, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal SKU from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct)
}

func CreateSKU(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	sku := new(models.SKU)

	if err := c.BodyParser(sku); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Basic validation
	if sku.SKU == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "SKU name is required",
		})
	}

	if sku.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Price must be greater than 0",
		})
	}

	// Set timestamps
	now := time.Now()
	sku.CreatedAt = now
	sku.UpdatedAt = now

	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Cannot fetch user ID - please log in",
		})
	}
	sku.UserID = userID
	sku.ID = uuid.New()

	result, count, err := supabaseClient.From("skus").Insert(sku, false, "", "", "").Execute() //I believe the other params are correct
	fmt.Println(string(result), count, err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save SKU to database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(sku)
}

func UpdateSKU(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	skuID := c.Params("id")
	sku := new(models.SKU)
	sid, err := uuid.Parse(skuID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid SKU ID",
		})
	}
	sku.ID = sid
	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Cannot fetch user ID - please log in",
		})
	}
	sku.UserID = userID

	if err := c.BodyParser(sku); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Basic validation
	if sku.SKU == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "SKU name is required",
		})
	}

	if sku.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Price must be greater than 0",
		})
	}

	// Set timestamps
	now := time.Now()
	sku.UpdatedAt = now

	//Save to database
	_, _, err = supabaseClient.From("skus").Update(sku, "", "").Eq("id", skuID).Execute()
	//fmt.Println(string(result), count, err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save SKU to database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(sku)
}

func DeleteSKU(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	skuID := c.Params("id")

	//Save to database
	_, _, err := supabaseClient.From("skus").Delete("", "").Eq("id", skuID).Execute()
	//fmt.Println(string(result), count, err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete SKU from database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "SKU deleted successfully",
	})
}
