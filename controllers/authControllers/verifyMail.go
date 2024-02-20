package authControllers

import (
	"encoding/json"
	"net/http"

	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

type Req struct {
	ConfirmationCode string `json:"confirmationCode"`
}

func VerifyMail(w http.ResponseWriter, r *http.Request) {
	var req Req
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error on verify email request", http.StatusBadRequest)
		return
	}

	user := &models.User{}
	database.DB.Where("verification_token = ?", req.ConfirmationCode).First(user)
	if user == nil {
		response := map[string]interface{}{
			"status":  "error",
			"message": "Invalid confirmation code",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if user.Verified == true {
		response := map[string]interface{}{
			"status":  "error",
			"message": "Email already verified",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if user.Verified == false {
		database.DB.Model(&models.User{}).Where("verification_token = ?", req.ConfirmationCode).Update("verified", true)
		response := map[string]interface{}{
			"status":  "success",
			"message": "Email verified",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}
