package projectsCRUD

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

func GetProject(c *fiber.Ctx) error {
	id := c.Params("id")

	//check if the user has access to this project

	var project models.Project

	database.DB.Where("id = ?", id).First(&project)

	return c.JSON(project)

}
 