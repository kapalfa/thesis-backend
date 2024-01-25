package authControllers

import (
	"encoding/json"
	"net/http"
	"github.com/kapalfa/go/models"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
	//"log"
	"fmt"
)

type EmailReq struct {
	Email string `json:"email"`
}
func LoginGithub(w http.ResponseWriter, r *http.Request) {
	var email EmailReq

	err := json.NewDecoder(r.Body).Decode(&email) 
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}

	fmt.Println("email:  ", email.Email)
	userModel, err := new(models.User), *new(error)
	userModel, err = getUserByEmail(email.Email)

	fmt.Println("userModel:  ", userModel)
	if userModel == nil {
	 	response := map[string]interface{}{
	 		"status": false,
	 		"message": "User not found",
	 	}
	 	json.NewEncoder(w).Encode(response)
	 	return
	} 

	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["id"] = userModel.Id
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token, err := accessToken.SignedString([]byte(config.Config("ACCESS_TOKEN_SECRET")))
	if err != nil {
	 	http.Error(w, err.Error(), http.StatusInternalServerError)
	 	return
	 }

	resfreshToken := jwt.New(jwt.SigningMethodHS256)
	claims = resfreshToken.Claims.(jwt.MapClaims)
	claims["id"] = userModel.Id
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	rt, err := resfreshToken.SignedString([]byte(config.Config("REFRESH_TOKEN_SECRET")))
	if err != nil {
	 	http.Error(w, err.Error(), http.StatusInternalServerError)
	 	return
	 }

	 userModel.RefreshToken = rt
	 database.DB.Save(&userModel)

	 cookie := http.Cookie{
	 	Name: "jwt",
	 	Value: rt,
	 	Expires: time.Now().Add(time.Hour),
	 	HttpOnly: true,
	 	Path: "/",
	 	SameSite: http.SameSiteNoneMode,
	 	Secure: true,
	 }

	 http.SetCookie(w, &cookie)

	 response := map[string]interface{}{
	 	"status": true,
	 	"message": "Logged in",
	 	"access_token": token,
	 	"cookie": cookie,
	 }

	 jsonResponse, err := json.Marshal(response)
	 if err != nil {
	 	http.Error(w, err.Error(), http.StatusInternalServerError)
	 	return
	 }
	 w.Header().Set("Content-Type", "application/json")
	 w.Write(jsonResponse)
}
