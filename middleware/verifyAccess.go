package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

func VerifyAccess() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("user")
		projectId := c.Params("id")
		var access models.Access
		err := database.DB.Model(&models.Access{}).Where("user_id = ? AND project_id = ?", userId, projectId).First(&access).Error
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Unauthenticated"}) 
		}
		return c.Next()
	}
}
