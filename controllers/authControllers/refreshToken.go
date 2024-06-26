package authControllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kapalfa/go/config"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		log.Print("No JWT cookie found")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "No JWT cookie found"})
		return
	}

	log.Println("cookie in refresh: ", cookie)
	var foundUser models.User
	database.DB.Model(&models.User{RefreshToken: cookie.Value}).First(&foundUser)
	if foundUser.Id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "User doesn't exist"})
		return
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config("REFRESH_TOKEN_SECRET")), nil
	})

	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Unauthorized user"})
		return
	}

	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["id"] = foundUser.Id
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix() // 15 minutes
	newToken, err := accessToken.SignedString([]byte(config.Config("ACCESS_TOKEN_SECRET")))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": newToken})
}
