package authControllers

import (
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"github.com/kapalfa/go/functions"
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userId, err := functions.GetUserId(auth)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	database.DB.Model(&models.User{}).Where("id=?", userId).Update("refresh_token", "")

	deletedCookie := http.Cookie{
		Name: "jwt",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path: "/",
		Secure: true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &deletedCookie)

	w.Write([]byte(`{"status": "success"}`))
}