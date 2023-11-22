package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"time"
)

func Logout(c *fiber.Ctx) error{
	//on client, also delete the access token when the logout button is clicked
	//delete the cookie from the client	

	cookie := c.Cookies("jwt")
	if cookie == "" {
		return c.Status(204).SendString("")
	}
	//is refresh token in DB ? if not just expire the cookie and return
	var foundUser models.User
	database.DB.Model(&models.User{RefreshToken : cookie}).First(&foundUser)
	if foundUser.Id == 0 {
		deletedCookie := fiber.Cookie{
			Name: "jwt",
			Value: "",
			Expires: time.Now().Add(-time.Hour),
			HTTPOnly: true,
			Path: "/",
			Secure: true,
			SameSite: "None",
		}
		c.Cookie(&deletedCookie)
		return c.Status(204).SendString("")
	}
	//if refresh token is in DB, delete it from DB and expire the cookie
	database.DB.Model(&models.User{}).Where("refresh_token = ?", cookie).Update("refresh_token", "")
	//create a cookie which expires in the past
	deletedCookie := fiber.Cookie{
		Name: "jwt",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path: "/",
		Secure: true,
		SameSite: "None",
	}
	c.Cookie(&deletedCookie)

	return c.JSON(fiber.Map{"status": "success"})
}