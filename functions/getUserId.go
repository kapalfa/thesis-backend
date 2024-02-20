package functions

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"os"
)

func GetUserId(authHeader string) (float64, error) {
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		return -1, fmt.Errorf("Invalid token")
	}
	token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return -1, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return -1, fmt.Errorf("Invalid token")
	}
	return claims["id"].(float64), nil
} 