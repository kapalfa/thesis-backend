package projectsCRUD

import (
	//"encoding/json"
	"net/http"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"github.com/gorilla/mux"
	"os"
)

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := database.DB.Where("project_id = ?", id).Delete(&models.Access{}).Error; err != nil { // delete all access to this project
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := database.DB.Where("id = ?", id).Delete(&models.Project{}).Error; err != nil { // delete project
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := os.RemoveAll("./uploads/" + id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, "Deleted project", http.StatusOK)

}