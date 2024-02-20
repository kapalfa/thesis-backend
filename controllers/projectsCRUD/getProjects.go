package projectsCRUD

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

type ProjectResponse struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

// get projects by user id
func GetProjects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["userid"]

	var projectIds []uint
	var projects []ProjectResponse
	database.DB.Model(&models.Access{}).Where("user_id = ? AND status = ?", id, "accepted").Pluck("project_id", &projectIds)
	if len(projectIds) == 0 {
		json.NewEncoder(w).Encode(projects)
		return
	}
	database.DB.Model(&models.Project{}).Where("id IN ?", projectIds).Select("Id", "Name", "Description", "Public").Find(&projects)
	json.NewEncoder(w).Encode(projects)
}
