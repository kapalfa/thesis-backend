package authControllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kapalfa/go/config"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func getUserByEmail(email string) (*models.User, error) {
	var user models.User
	db := database.DB
	if err := db.Where(&models.User{Email: email}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, errors.New("user not verified")
	}
	if !user.Verified {
		return nil, errors.New("user not verified")
	}
	return &user, nil
}
func CheckPasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	//request body
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type UserData struct {
		Id           uint   `json:"id"`
		Email        string `json:"email"`
		Password     []byte `json:"password"`
		RefreshToken string `json:"refresh_token"`
	}
	var userData UserData
	email := request.Email
	pass := request.Password
	userModel, err := new(models.User), *new(error)
	userModel, err = getUserByEmail(email)

	if err != nil {
		if err == gorm.ErrRecordNotFound { // email not found
			response := map[string]interface{}{
				"message": "Email not found",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		} else if err.Error() == "user not verified" { // email not verified
			response := map[string]interface{}{
				"message": "User not verified",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		} else if userModel == nil { // other error
			response := map[string]interface{}{
				"message": "User not found",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		userData = UserData{
			Id:       userModel.Id,
			Email:    userModel.Email,
			Password: []byte(userModel.Password),
		}
	}

	if !CheckPasswordHash(pass, userData.Password) { // incorrect password
		response := map[string]interface{}{
			"message": "Invalid password",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	//create access token
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["id"] = userData.Id
	claims["exp"] = time.Now().Add(15 * time.Minute).Unix() // 15 minutes
	token, err := accessToken.SignedString([]byte(config.Config("ACCESS_TOKEN_SECRET")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//create refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claims = refreshToken.Claims.(jwt.MapClaims)
	claims["id"] = userData.Id
	claims["exp"] = time.Now().Add(time.Hour).Unix() // 1 hour
	rt, err := refreshToken.SignedString([]byte(config.Config("REFRESH_TOKEN_SECRET")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//update user entry to store refresh token
	userModel.RefreshToken = rt
	database.DB.Save(&userModel)

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    rt,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	response := map[string]interface{}{
		"status":       "success",
		"message":      "Logged in",
		"access_token": token,
		"cookie":       cookie,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
