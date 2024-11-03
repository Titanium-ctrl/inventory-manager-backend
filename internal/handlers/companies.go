package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
	"ucrs.com/inventory-manager/backend/internal/models"
)

func CreateCompany(c *fiber.Ctx) error {
	company := new(models.Company)

	if err := c.BodyParser(company); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)

	if company.ID == uuid.Nil || company.Name == "" || company.Industry == "" || company.Owner == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()

	//Insert user
	_, _, err := supabaseClient.From("companies").Insert(company, false, "", "", "").Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot save company to database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Company created successfully",
	})
}

/*
func GetCompanies(c *fiber.Ctx) error {


	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Company created successfully",
	})
}
*/

func GetCompany(c *fiber.Ctx) error {
	companyid := c.Params("id")
	_, err := uuid.Parse(companyid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid company ID",
		})
	}
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	//Fetch user from supabase
	user, _, err := supabaseClient.From("companies").Select("*", "", false).Eq("id", companyid).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch company from database",
		})
	}
	respStruct := []struct {
		models.Company
	}{}
	err = json.Unmarshal(user, &respStruct)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal company from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(respStruct[0])
}

func UpdateCompany(c *fiber.Ctx) error {
	companyid := c.Params("id")
	if companyid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Company ID is required",
		})
	}
	var company models.Company
	err := c.BodyParser(&company)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot unmarshal company from request body",
		})
	}
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)
	cid, err := uuid.Parse(companyid)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot parse company ID",
		})
	}
	company.ID = cid
	_, _, err = supabaseClient.From("companies").Update(company, "", "").Eq("id", companyid).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update company in database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Company updated successfully",
	})
}

func DeleteCompany(c *fiber.Ctx) error {
	companyid := c.Params("id")
	if companyid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Company ID is required",
		})
	}
	_, err := uuid.Parse(companyid)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot parse company ID",
		})
	}
	supabaseClient := c.Locals("supabaseClient").(*supabase.Client)

	_, _, err = supabaseClient.From("companies").Delete("", "").Eq("id", companyid).Execute()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete company from database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Company deleted successfully",
	})
}
