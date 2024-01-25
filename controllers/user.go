package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

// func User(c *fiber.Ctx) error {
// 	userId := c.Locals("user")
// 	user := models.User{}
// 	database.DB.First(&user, userId)
	
// 	return c.JSON(fiber.Map{"status":"success", "message": "Access to user profile", "user info": user})
// }

func User(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["user"]
	user := User{}
	DB.First(&user, userId)

	response := map[string]interface{}{
		"status": "success",
		"message": "Access to user profile"
		"user info": user
	}
	json.NewEncoder(w).Encode(response)
}