package projectsCRUD

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

func GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var project models.Project

	database.DB.Where("id = ?", id).First(&project)

	json.NewEncoder(w).Encode(project)
}
