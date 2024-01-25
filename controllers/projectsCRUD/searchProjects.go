package projectsCRUD

import (
	"encoding/json"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"github.com/gorilla/mux"
	"net/http"
)

func SearchProjects(w http.ResponseWriter, r *http.Request) {
	debouncedSearch := mux.Vars(r)["projectName"]

	var projects []models.Project
	database.DB.Where("name LIKE ? AND public = ?", "%"+debouncedSearch+"%", true).Find(&projects)
	json.NewEncoder(w).Encode(projects)
}