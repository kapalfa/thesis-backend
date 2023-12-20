package authControllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/kapalfa/go/config"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

func getUserByEmail(email string) (*models.User, error) {
	var user models.User
	db := database.DB
	if err := db.Where(&models.User{Email: email}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func CheckPasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(c *fiber.Ctx) error {
	//request body 
	type LoginRequest struct {
		Email 		string	`json:"email"`
		Password 	string 	`json:"password"`
	}
	request := new(LoginRequest)
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	type UserData struct {
		Id 			uint 	`json:"id"`
		Email 		string 	`json:"email"`
		Password 	string 	`json:"password"`
		RefreshToken string `json:"refresh_token"`
	}
	var userData UserData
	email := request.Email
	pass := request.Password
	userModel, err := new(models.User), *new(error)
	userModel, err = getUserByEmail(email)
	
	if userModel == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "User not found", "data": err})
	} else {
		userData = UserData{
			Id: userModel.Id,
			Email: userModel.Email,
			Password: userModel.Password,
		}
	}
	
	if !CheckPasswordHash(pass, []byte(userData.Password)) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Incorrect password", "data": nil})
	}
	
	//create access token
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["id"] = userData.Id
	claims["exp"] = time.Now().Add(time.Minute*15).Unix() // 15 minutes
	token, err := accessToken.SignedString([]byte(config.Config("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	//create refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claims = refreshToken.Claims.(jwt.MapClaims)
	claims["id"] = userData.Id
	claims["exp"] = time.Now().Add(time.Hour).Unix() // 1 hour
	rt, err := refreshToken.SignedString([]byte(config.Config("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	//update user entry to store refresh token
	userModel.RefreshToken = rt
	database.DB.Save(&userModel)
	//send cookie to user
	cookie := &fiber.Cookie{
	 	Name : "jwt",
	 	Value: rt,
	 	Expires: time.Now().Add(time.Hour),
	 	HTTPOnly: true,
		Path: "/",
		SameSite: "None",
		Secure: true,
	}
	c.Cookie(cookie)
	return c.JSON(fiber.Map{"status": "success", "message": "Logged in", "access_token": token, "cookie": cookie})
}