package middleware

import (
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"github.com/gorilla/mux"
	"net/http"
)

func VerifyAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user")
		projectId := mux.Vars(r)["id"]
		var access models.Access
		err := database.DB.Model(&models.Access{}).Where("user_id = ? AND project_id = ?", userId, projectId).First(&access).Error
		if err != nil {
			http.Error(w, "Unauthenticated", http.StatusUnauthorized) 
			return
		}
		next.ServeHTTP(w, r)
	})
}
