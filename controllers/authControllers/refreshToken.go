 package authControllers

 import (
 	"github.com/gofiber/fiber/v2"
 	"github.com/kapalfa/go/config"
 	"github.com/kapalfa/go/database"
 	"github.com/kapalfa/go/models"
 	"github.com/golang-jwt/jwt/v5"
 	"time"
 )

 func HandleRefreshToken(c *fiber.Ctx) error { 
	
	cookie := c.Cookies("jwt")
	if cookie == "" {	
		return c.JSON(fiber.Map{"token": "error"})
	}

	var foundUser models.User
	database.DB.Model(&models.User{RefreshToken : cookie}).First(&foundUser)	
	if foundUser.Id == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Unauthorized", "data": nil})
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config("REFRESH_TOKEN_SECRET")), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Unauthorized", "data": nil})
	}
	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Unauthorized", "data": nil})
	}

	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["id"] = foundUser.Id
	claims["exp"] = time.Now().Add(time.Minute*15).Unix() // 15 minutes
	newToken, err := accessToken.SignedString([]byte(config.Config("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": newToken})
}