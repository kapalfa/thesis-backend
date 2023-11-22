package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"os"
)

//verify access token and output the claims
func VerifyJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Unauthenticated"})
		}
		tokenString := strings.Split(authHeader, " ")[1]
		mySigningKey := []byte(os.Getenv("ACCESS_TOKEN_SECRET"))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return mySigningKey, nil
		})
		if err != nil  || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Unauthenticated"})
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Locals("user", claims["id"])
		} else {
			fmt.Println(err)
		}
	return c.Next()	}
} 