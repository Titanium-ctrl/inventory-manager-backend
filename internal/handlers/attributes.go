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

// CreateAttribute
func CreateAttribute(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	attribute := new(models.Attribute)

	if err := c.BodyParser(attribute); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	//Attribute validation - attribute should have a name and user_id
	if attribute.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Attribute name is required",
		})
	}

	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Cannot fetch user ID - please log in",
		})
	}
	attribute.UserID = userID
	attribute.ID = uuid.New()

	// Set timestamps
	now := time.Now()
	attribute.CreatedAt = now

	//Save to database
	_, _, err = supabaseClient.From("attributes").Insert(attribute, false, "", "", "").Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save attribute to database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Attribute created successfully",
		"attribute": attribute,
	})
}

// GetAttributes
func GetAttributes(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	attributes, _, err := supabaseClient.From("attributes").Select("*", "", false).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch attributes from database",
		})
	}
	respStruct := []struct {
		models.Attribute
	}{}
	err = json.Unmarshal(attributes, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal attribute from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct)
}

// GetAttribute
func GetAttribute(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	attributeID := c.Params("id")
	attribute, _, err := supabaseClient.From("attributes").Select("*", "", false).Eq("id", attributeID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch attribute from database",
		})
	}
	respStruct := []struct {
		models.Attribute
	}{}

	err = json.Unmarshal(attribute, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal attribute from database",
		})
	}

	if len(respStruct) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Attribute not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct[0])
}

// UpdateAttribute
func UpdateAttribute(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	attributeID := c.Params("id")
	attribute := new(models.Attribute)
	aid, err := uuid.Parse(attributeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid attribute ID",
		})
	}
	attribute.ID = aid
	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Cannot fetch user ID - please log in",
		})
	}
	attribute.UserID = userID

	if err := c.BodyParser(attribute); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	//Attribute validation - attribute should have a name and user_id
	if attribute.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Attribute name is required",
		})
	}

	//Save to database
	_, _, err = supabaseClient.From("attributes").Update(attribute, "", "").Eq("id", attributeID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save attribute to database",
		})
	}

	//If invalid user is making request, it doesn't apply the change, but the below response is still returned
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Attribute updated successfully",
	})
}

// DeleteAttribute
func DeleteAttribute(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	attributeID := c.Params("id")

	//Save to database
	_, _, err := supabaseClient.From("attributes").Delete("", "").Eq("id", attributeID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete attribute from database",
		})
	}

	//Positive result shows even if RLS policy blocks it
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Attribute deleted successfully",
	})
}
