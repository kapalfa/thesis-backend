package controllers

import (
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"net/http"
	"github.com/gorilla/mux"
	"strings"
)

func VerifyAccess(w http.ResponseWriter, r *http.Request){
	projectId := r.URL.Query().Get("projectid")
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	userId := getUserId(token)
	var access models.Access
	err := database.DB.Model(&models.Access{}).Where("user_id = ? AND project_id = ?", userId, projectId).First(&access).Error
	if err != nil {
		http.Error(w, "Unauthorized", 401)
		return
	}
}