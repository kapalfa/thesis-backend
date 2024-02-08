package projectsCRUD

import (
	"encoding/json"
	"net/http"	

	"github.com/gorilla/mux"
	"github.com/kapalfa/go/models"
	"github.com/kapalfa/go/database"
)

func GetCollaborators(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId := vars["id"]
	userId := vars["userid"]

	// get collaborators from database
	var userids []uint
	database.DB.Select("user_id").Where("project_id = ? AND user_id != ?", projectId, userId).Find(&models.Access{}).Pluck("user_id", &userids)
	
	var emails []string
	database.DB.Select("email").Where("id IN (?)", userids).Find(&models.User{}).Pluck("email", &emails)
	json.NewEncoder(w).Encode(emails)
}