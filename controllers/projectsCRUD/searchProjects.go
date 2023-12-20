package projectsCRUD

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)
type SearchRequest struct {
	Search string `json:"search"`
}

func SearchProjects(c *fiber.Ctx) error {
	debouncedSearch := c.Params("projectName")
	//var request SearchRequest
	//if err := c.BodyParser(&request); err != nil {
	//	return c.Status(400).SendString(err.Error())
	//}

	fmt.Println("search for : ", debouncedSearch)
	var projects []models.Project
	database.DB.Where("name LIKE ? AND public = ?", "%"+debouncedSearch+"%", true).Find(&projects)
	return c.JSON(projects)
}