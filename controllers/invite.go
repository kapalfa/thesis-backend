package controllers 

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

func Invite(c *fiber.Ctx) error {
	var input map[string]interface{}
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Couldn't invite user", "data": err})
	}

	access := &models.Access{
		UserId:    uint(input["user_id"].(float64)),
		ProjectId: uint(input["project_id"].(float64)),
	}

	if err := database.DB.Create(access).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't invite user", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Invited user", "data": access})
}