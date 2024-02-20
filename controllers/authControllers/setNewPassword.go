package authControllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"golang.org/x/crypto/bcrypt"
)

type NewPasswordReq struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func SetNewPassword(w http.ResponseWriter, r *http.Request) {
	var req NewPasswordReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error on set new password request", http.StatusBadRequest)
		return
	}

	var user models.User
	log.Print("new password : ", req.Password)
	log.Print("reset token : ", req.Token)
	database.DB.Where("reset_token = ?", req.Token).First(&user)
	if user.Id == 0 {
		response := map[string]interface{}{
			"status":  "error",
			"message": "Invalid token",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if user.ResetTokenExpires.Before(time.Now()) {
		response := map[string]interface{}{
			"status":  "error",
			"message": "Token expired",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		http.Error(w, "Couldn't hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	log.Print("new password : ", user.Password)
	database.DB.Save(&user)
	var myuser models.User
	res := database.DB.Where("email = ?", user.Email).First(&myuser)
	if res.Error != nil {
		http.Error(w, "Error on set new password request", http.StatusInternalServerError)
		return
	}
	log.Print("user : ", myuser)
	response := map[string]interface{}{
		"status":  "success",
		"message": "Password updated",
	}
	json.NewEncoder(w).Encode(response)
}
