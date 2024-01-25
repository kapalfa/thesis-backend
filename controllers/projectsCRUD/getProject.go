package projectsCRUD

import (
	"encoding/json"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
)

func GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var project models.Project

	fmt.Println("get project with id : ", id)

	database.DB.Where("id = ?", id).First(&project)

	fmt.Println("project : ", project)
	json.NewEncoder(w).Encode(project)

}
 