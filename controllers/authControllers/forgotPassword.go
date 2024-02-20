package authControllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/functions"
	"github.com/kapalfa/go/models"
	"github.com/kapalfa/go/utils"
)

type Request struct {
	Email string `json:"email"`
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error on forgot password request", http.StatusBadRequest)
		return
	}
	user := &models.User{}
	database.DB.Where("email = ?", req.Email).First(user)
	if user.Id == 0 {
		response := map[string]interface{}{
			"status":  "error",
			"message": "User not found",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	//send email to the user with token
	resetToken := utils.GenerateRandomString(16) // create token
	resetTokenExpires := time.Now().Add(time.Hour * 1)
	database.DB.Model(&models.User{}).Where("email = ?", req.Email).Updates(map[string]interface{}{
		"reset_token":         resetToken,
		"reset_token_expires": resetTokenExpires,
	})
	functions.SendEmail(req.Email, resetToken, "reset-password")
	response := map[string]interface{}{
		"status":  "success",
		"message": "Email sent to reset password",
	}
	json.NewEncoder(w).Encode(response)
}
