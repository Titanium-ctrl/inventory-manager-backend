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

func convertWarehouseForDB(warehouse *models.Warehouse) *models.WarehouseDatabase {
	return &models.WarehouseDatabase{
		ID:           warehouse.ID,
		UserID:       warehouse.UserID,
		Name:         warehouse.Name,
		AddressLine1: warehouse.Address.AddressLine1,
		AddressLine2: warehouse.Address.AddressLine2,
		TownCity:     warehouse.Address.TownCity,
		StateCounty:  warehouse.Address.StateCounty,
		PostZipCode:  warehouse.Address.PostZipCode,
		Country:      warehouse.Address.Country,
		Latitude:     warehouse.Latitude,
		Longitude:    warehouse.Longitude,
		CreatedAt:    warehouse.CreatedAt,
		UpdatedAt:    warehouse.UpdatedAt,
	}
}

func convertWarehouseForJSON(warehouse *models.WarehouseDatabase) *models.Warehouse {
	res := &models.Warehouse{
		ID:        warehouse.ID,
		UserID:    warehouse.UserID,
		Name:      warehouse.Name,
		Latitude:  warehouse.Latitude,
		Longitude: warehouse.Longitude,
		CreatedAt: warehouse.CreatedAt,
		UpdatedAt: warehouse.UpdatedAt,
	}
	res.Address.AddressLine1 = warehouse.AddressLine1
	res.Address.AddressLine2 = warehouse.AddressLine2
	res.Address.TownCity = warehouse.TownCity
	res.Address.StateCounty = warehouse.StateCounty
	res.Address.PostZipCode = warehouse.PostZipCode
	res.Address.Country = warehouse.Country
	return res
}

func CreateWarehouse(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	inputwarehouse := new(models.Warehouse)

	if err := c.BodyParser(inputwarehouse); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if inputwarehouse.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Warehouse name is required",
		})
	}

	warehouse := convertWarehouseForDB(inputwarehouse)

	warehouse.ID = uuid.New()

	// Set timestamps
	now := time.Now()
	warehouse.CreatedAt = now
	warehouse.UpdatedAt = now

	//Fetch userID and apply to warehouse
	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch user ID from database",
		})
	}
	warehouse.UserID = userID

	//Save to database
	_, _, err = supabaseClient.From("warehouses").Insert(warehouse, false, "", "", "").Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save warehouse to database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(warehouse)
}

func UpdateWarehouse(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	warehouseID := c.Params("id")
	inputwarehouse := new(models.Warehouse)
	wid, err := uuid.Parse(warehouseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid warehouse ID",
		})
	}

	if err := c.BodyParser(inputwarehouse); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Basic validation
	if inputwarehouse.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Warehouse name is required",
		})
	}

	//addressline1, towncity, postzipcode, country
	if inputwarehouse.Address.AddressLine1 == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Address line 1 is required",
		})
	}

	if inputwarehouse.Address.TownCity == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Town city is required",
		})
	}

	if inputwarehouse.Address.PostZipCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Post zip code is required",
		})
	}

	if inputwarehouse.Address.Country == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Country is required",
		})
	}

	warehouse := convertWarehouseForDB(inputwarehouse)

	warehouse.ID = wid

	userid, err := database.FetchUserID(supabaseClient)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch user ID from database",
		})
	}
	warehouse.UserID = userid

	// Set timestamps
	now := time.Now()
	warehouse.UpdatedAt = now

	//Save to database
	_, _, err = supabaseClient.From("warehouses").Update(warehouse, "", "").Eq("id", warehouseID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save warehouse to database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(warehouse)
}

func GetWarehouses(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	warehouses, _, err := supabaseClient.From("warehouses").Select("*", "", false).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch warehouses from database",
		})
	}
	respStruct := []models.WarehouseDatabase{}
	err = json.Unmarshal(warehouses, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal warehouse from database",
		})
	}

	// Create a slice to hold the converted warehouses
	respwarehouses := make([]*models.Warehouse, len(respStruct))

	// Iterate through the dbSlice and convert each *WarehouseDatabase to *Warehouse
	for i, db := range respStruct {
		respwarehouses[i] = convertWarehouseForJSON(&db)
	}

	return c.Status(fiber.StatusOK).JSON(respwarehouses)
}

func GetWarehouse(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	warehouseID := c.Params("id")
	warehouse, _, err := supabaseClient.From("warehouses").Select("*", "", false).Eq("id", warehouseID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch warehouse from database",
		})
	}
	respStruct := []models.WarehouseDatabase{}
	err = json.Unmarshal(warehouse, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal warehouse from database",
		})
	}

	resp := convertWarehouseForJSON(&respStruct[0])

	return c.Status(fiber.StatusOK).JSON(resp)
}

func DeleteWarehouse(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	warehouseID := c.Params("id")

	//Save to database
	_, _, err := supabaseClient.From("warehouses").Delete("", "").Eq("id", warehouseID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete warehouse from database",
		})
	}

	//Positive result shows even if RLS policy blocks it
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Warehouse deleted successfully",
	})
}
