package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
	"ucrs.com/inventory-manager/backend/internal/database"
	"ucrs.com/inventory-manager/backend/internal/models"
)

func UpdateSKUAttribute(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	skuID := c.Params("skuid")

	SKUAttr := new(models.SKUAttributes)
	if err := c.BodyParser(SKUAttr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if SKUAttr.AttributeID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Attribute ID is required",
		})
	}

	if SKUAttr.AttributeValue == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Attribute value is required",
		})
	}

	SKUID, err := uuid.Parse(skuID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sku ID",
		})
	}
	SKUAttr.SkuID = SKUID

	//Get UserID
	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch user ID from database",
		})
	}
	SKUAttr.UserID = userID

	_, _, err = supabaseClient.From("sku_attributes").Upsert(SKUAttr, "sku_id, attribute_id", "", "").Eq("sku_id", skuID).Eq("attribute_id", SKUAttr.AttributeID.String()).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update inventory in database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "SKU Attribute updated successfully",
	})
}

func GetSKUAttributes(c *fiber.Ctx) error {
	skuid := c.Params("skuid")
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)

	attributes, _, err := supabaseClient.From("sku_attributes").Select("*", "", false).Eq("sku_id", skuid).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch inventory from database",
		})
	}
	respStruct := []models.SKUAttributes{}
	err = json.Unmarshal(attributes, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal inventory from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "SKU Attributes retrieved successfully",
		"response": respStruct,
	})
}

func GetSKUAttribute(c *fiber.Ctx) error {
	skuid := c.Params("skuid")
	attributeid := c.Params("id")
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)

	attributes, _, err := supabaseClient.From("sku_attributes").Select("*", "", false).Eq("sku_id", skuid).Eq("attribute_id", attributeid).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch inventory from database",
		})
	}
	respStruct := []models.SKUAttributes{}
	err = json.Unmarshal(attributes, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal inventory from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "SKU Attributes retrieved successfully",
		"response": respStruct,
	})
}

func DeleteSKUAttribute(c *fiber.Ctx) error {
	skuid := c.Params("skuid")
	attributeid := c.Params("id")
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)

	_, _, err := supabaseClient.From("sku_attributes").Delete("", "").Eq("sku_id", skuid).Eq("attribute_id", attributeid).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete inventory from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "SKU Attribute deleted successfully",
	})
}
