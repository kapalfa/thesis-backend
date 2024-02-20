package projectsCRUD

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

func GetCollaborators(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId := vars["id"]
	userId := vars["userid"]
	// get collaborators from database
	var userids []uint
	database.DB.Select("user_id").Where("project_id = ? AND user_id != ? AND status = ?", projectId, userId, "accepted").Find(&models.Access{}).Pluck("user_id", &userids)
	if len(userids) == 0 {
		json.NewEncoder(w).Encode([]string{})
		log.Println("No collaborators found")
		return
	}
	var emails []string
	database.DB.Select("email").Where("id IN (?)", userids).Find(&models.User{}).Pluck("email", &emails)
	json.NewEncoder(w).Encode(emails)
}
