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

func UpdateInventory(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	locationID := c.Params("locationid")
	skuID := c.Params("skuid")
	inventory := new(models.Inventory)
	if err := c.BodyParser(inventory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	locID, err := uuid.Parse(locationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid location ID",
		})
	}
	inventory.LocationID = locID

	inventory.SkuID, err = uuid.Parse(skuID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid sku ID",
		})
	}

	// Set timestamps
	currentTime := time.Now()
	inventory.UpdatedAt = currentTime

	//Get UserID
	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch user ID from database",
		})
	}
	inventory.UserID = userID

	_, _, err = supabaseClient.From("inventory").Upsert(inventory, "sku_id, location_id", "", "").Eq("location_id", locationID).Eq("sku_id", skuID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update inventory in database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(inventory)
}

// FUNCTION WILL RETURN TOO MUCH, OR JUST UP TO SUPABASE LIMITS = ADD PAGINATION FROM RESUABLE COMPONENT!
func GetInventory(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	locationid := c.Params("locationid")

	respStruct := []models.Inventory{}

	if locationid == "" {
		inventory, _, err := supabaseClient.From("inventory").Select("*", "", false).Execute()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot fetch inventory from database",
			})
		}
		err = json.Unmarshal(inventory, &respStruct)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot unmarshal inventory from database",
			})
		}

		/*
			// Create a slice to hold the converted inventory
			respinventory := make([]*models.Inventory, len(respStruct))

			// Iterate through the dbSlice and convert each *InventoryDatabase to *Inventory
			for i, db := range respStruct {
				respinventory[i] = convertInventoryForJSON(&db)
			}
		*/
	} else {
		inventory, _, err := supabaseClient.From("inventory").Select("*", "", false).Eq("location_id", locationid).Execute()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot fetch inventory from database",
			})
		}
		err = json.Unmarshal(inventory, &respStruct)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot unmarshal inventory from database",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(respStruct)
}

func GetSpecificInventory(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	locationid := c.Params("locationid")
	skuid := c.Params("skuid")

	inventory, _, err := supabaseClient.From("inventory").Select("*", "", false).Eq("location_id", locationid).Eq("sku_id", skuid).Execute()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch inventory from database",
		})
	}
	respStruct := []models.Inventory{}
	err = json.Unmarshal(inventory, &respStruct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal inventory from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct[0])
}

func GetInventoryForSKU(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	skuID := c.Params("skuid")

	inventory, _, err := supabaseClient.From("inventory").Select("*", "", false).Eq("sku_id", skuID).Execute()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch inventory from database",
		})
	}
	respStruct := []models.Inventory{}
	err = json.Unmarshal(inventory, &respStruct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal inventory from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct)
}

func DeleteInventory(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	locationID := c.Params("locationid")

	_, err := uuid.Parse(locationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid location ID",
		})
	}

	_, _, err = supabaseClient.From("inventory").Delete("", "").Eq("location_id", locationID).Execute()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete inventory from database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Inventory deleted successfully",
	})
}
