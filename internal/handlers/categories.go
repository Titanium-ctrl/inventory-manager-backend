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

func CreateCategory(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	category := new(models.Category)

	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if category.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category name is required",
		})
	}

	category.ID = uuid.New()

	//Fetch userID and apply to category
	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch user ID from database",
		})
	}
	category.UserID = userID

	//Save to database
	_, _, err = supabaseClient.From("categories").Insert(category, false, "", "", "").Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save category to database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

func UpdateCategory(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	categoryID := c.Params("id")
	category := new(models.Category)
	cid, err := uuid.Parse(categoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	category.ID = cid

	// Basic validation
	if category.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category name is required",
		})
	}

	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch user ID from database",
		})
	}
	category.UserID = userID

	//Save to database
	_, _, err = supabaseClient.From("categories").Update(category, "", "").Eq("id", categoryID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save category to database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(category)
}

func GetCategories(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	categories, _, err := supabaseClient.From("categories").Select("*", "", false).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch categories from database",
		})
	}
	respStruct := []struct {
		models.Category
	}{}
	err = json.Unmarshal(categories, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal category from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct)
}

// GetCategoriesByParentID
func GetCategoriesByParentID(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	parentID := c.Params("id")
	categories, _, err := supabaseClient.From("categories").Select("*", "", false).Eq("parent_id", parentID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch categories from database",
		})
	}
	respStruct := []struct {
		models.Category
	}{}
	err = json.Unmarshal(categories, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal category from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct)
}

func GetCategory(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	categoryID := c.Params("id")
	category, _, err := supabaseClient.From("categories").Select("*", "", false).Eq("id", categoryID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch category from database",
		})
	}
	respStruct := []struct {
		models.Category
	}{}
	err = json.Unmarshal(category, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal category from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct[0])
}

func DeleteCategory(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	categoryID := c.Params("id")

	//Save to database
	_, _, err := supabaseClient.From("categories").Delete("", "").Eq("id", categoryID).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete category from database",
		})
	}

	//Positive result shows even if RLS policy blocks it
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}
