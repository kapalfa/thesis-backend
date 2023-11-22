package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	
	c.Set("Content-Type", "application/json")
	type NewUser struct {
		Name 		string `json:"name"`
		Email 		string `json:"email"`
	}
	
	user :=  &models.User{}

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on register request", "data": err})
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})
	}
	
	user.Password = string(password)
	if err := database.DB.Create(user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}
	newUser := NewUser{
		Name: user.Name,
		Email: user.Email,
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": newUser})
}