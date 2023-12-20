package projectsCRUD

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

type ProjectResponse struct {
	Id uint `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Public bool `json:"public"`
}

//get projects of user based on the user id
func GetProjects(c *fiber.Ctx) error {
	id := c.Params("userid")
	var projectIds []uint
	var projects []ProjectResponse
	database.DB.Model(&models.Access{}).Where("user_id = ?", id).Pluck("project_id", &projectIds)

	database.DB.Model(&models.Project{}).Where("id IN ?", projectIds).Select("Id", "Name", "Description", "Public").Find(&projects)

	return c.JSON(projects)
}