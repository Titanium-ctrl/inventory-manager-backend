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

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)

	if user.FirstName == "" || user.LastName == "" || user.CompanyID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	id, err := database.FetchUserID(supabaseClient)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch user ID",
		})
	}
	user.ID = id
	user.UpdatedAt = time.Now()

	//Insert user
	_, _, err = supabaseClient.From("users").Insert(user, false, "", "", "").Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save user to database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user":    user,
	})

}

func GetUsersFromCompanyID(c *fiber.Ctx) error {
	// CHECK USER PERMISSIONS FIRST!!!!!!!!!!!

	companyid := c.Params("companyid")
	_, err := uuid.Parse(companyid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid company ID",
		})
	}
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)

	users, _, err := supabaseClient.From("users").Select("*", "", false).Eq("company_id", companyid).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch users from database",
		})
	}

	respStruct := []struct {
		models.User
	}{}
	err = json.Unmarshal(users, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal user from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct)
}

func GetUser(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)

	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch user ID",
		})
	}

	//Fetch user from supabase
	user, _, err := supabaseClient.From("users").Select("*", "", false).Eq("id", userID.String()).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch user from database",
		})
	}
	respStruct := []struct {
		models.User
	}{}
	err = json.Unmarshal(user, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal user from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct[0])
}

func UpdateUser(c *fiber.Ctx) error {

	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch user ID",
		})
	}

	var user models.User
	err = c.BodyParser(&user)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal user from request body",
		})
	}
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot parse user ID",
		})
	}
	user.ID = userID
	_, _, err = supabaseClient.From("users").Update(user, "", "").Eq("id", userID.String()).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update user in database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	userID, err := database.FetchUserID(supabaseClient)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch user ID",
		})
	}

	_, _, err = supabaseClient.From("users").Delete("", "").Eq("id", userID.String()).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete user from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
