package authControllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/functions"
	"github.com/kapalfa/go/models"
	"github.com/kapalfa/go/utils"
)

type NewUser struct {
	Email string `json:"email"`
}

func Register(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, "Error on register request", http.StatusBadRequest)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		http.Error(w, "Couldn't hash password", http.StatusInternalServerError)
		return
	}

	user.Password = string(password)
	verificationToken := utils.GenerateRandomString(16) // create token for email verification
	user.VerificationToken = verificationToken
	if err := database.DB.Create(user).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { //postgresql code for duplicate entry
				response := map[string]interface{}{
					"status":  "error",
					"message": "User already exists",
				}
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}
	//send email to the user with token
	functions.SendEmail(user.Email, verificationToken, "verify-email")

	newUser := NewUser{
		Email: user.Email,
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Created user",
		"data":    newUser,
	}
	json.NewEncoder(w).Encode(response)
}
