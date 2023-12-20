package projectsCRUD

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"os"
	"strconv"
)

func CreateProject(c *fiber.Ctx) error {
	var input map[string]interface{}
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Couldn't create project", "data": err})
	}
	
	project :=  &models.Project{
		Name: input["name"].(string),
		Description: input["description"].(string),
		Public: input["public"].(bool),
	}

	if err := database.DB.Create(project).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create project", "data": err})
	}

	if err := os.MkdirAll("./uploads/" + strconv.Itoa(int(project.Id)), 0755); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create project directory", "data": err})
	}

	type NewProject struct {
		Id 			uint   `json:"id"`
		Name 		string `json:"name"`
		Description string `json:"description"`
		Public 		bool   `json:"public"`
	}

	newProject := NewProject{
		Id:          project.Id,
		Name:        project.Name,
		Description: project.Description,
		Public:      project.Public,
	}

	access := &models.Access{
		UserId:    uint(input["user_id"].(float64)),
		ProjectId: project.Id,
	}

	if err := database.DB.Create(access).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create access on project", "data": err})
	}

	
	return c.JSON(fiber.Map{"status": "success", "message": "Created project", "data": newProject})
}

