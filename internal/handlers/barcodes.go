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

func CreateBarcode(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	barcode := new(models.Barcode)

	if err := c.BodyParser(barcode); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if barcode.BarcodeName == "" {
		//Set default barcode name
		barcode.BarcodeName = "Product Barcode"
	}

	if barcode.BarcodeValue == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Barcode value is required",
		})
	}

	if barcode.SkuID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "SKU ID is required",
		})
	}

	//Fetcch userID and apply to barcode
	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch user ID from database",
		})
	}
	barcode.UserID = userID

	// Set timestamps
	now := time.Now()
	barcode.CreatedAt = now
	barcode.UpdatedAt = now

	barcode.ID = uuid.New()

	//Save to database
	_, _, err = supabaseClient.From("barcodes").Insert(barcode, false, "", "", "").Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save barcode to database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Barcode created successfully",
		"barcode": barcode,
	})
}

func UpdateBarcode(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	barcodeID := c.Params("id")
	barcode := new(models.Barcode)
	bid, err := uuid.Parse(barcodeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid barcode ID",
		})
	}

	if err := c.BodyParser(barcode); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	barcode.ID = bid

	// Basic validation
	if barcode.BarcodeName == "" {
		//Set default barcode name
		barcode.BarcodeName = "Product Barcode"
	}

	if barcode.BarcodeValue == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Barcode value is required",
		})
	}

	if barcode.SkuID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "SKU ID is required",
		})
	}

	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch user ID from database",
		})
	}
	barcode.UserID = userID

	// Set timestamps
	now := time.Now()
	barcode.UpdatedAt = now

	//Save to database
	_, _, err = supabaseClient.From("barcodes").Update(barcode, "", "").Eq("id", barcodeID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save updated barcode to database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Barcode updated successfully",
	})
}

func GetBarcodes(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	barcodes, _, err := supabaseClient.From("barcodes").Select("*", "", false).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch barcodes from database",
		})
	}
	respStruct := []struct {
		models.Barcode
	}{}
	err = json.Unmarshal(barcodes, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal barcode from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct)
}

// GetBarcodesBySKUID
func GetBarcodesBySKUID(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	skuID := c.Params("id")
	barcodes, _, err := supabaseClient.From("barcodes").Select("*", "", false).Eq("sku_id", skuID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch barcodes from database",
		})
	}
	respStruct := []struct {
		models.Barcode
	}{}
	err = json.Unmarshal(barcodes, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal barcode from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct)
}

// GetBarcode
func GetBarcode(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	barcodeID := c.Params("id")
	barcode, _, err := supabaseClient.From("barcodes").Select("*", "", false).Eq("id", barcodeID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch barcode from database",
		})
	}
	respStruct := []struct {
		models.Barcode
	}{}
	err = json.Unmarshal(barcode, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal barcode from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct[0])
}

// DeleteBarcode
func DeleteBarcode(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	barcodeID := c.Params("id")

	//Save to database
	_, _, err := supabaseClient.From("barcodes").Delete("", "").Eq("id", barcodeID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete barcode from database",
		})
	}

	//Positive result shows even if RLS policy blocks it
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Barcode deleted successfully",
	})
}
