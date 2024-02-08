package projectsCRUD

import (
	"encoding/json"
	"net/http"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"os"
	"strconv"
)

func CreateProject(w http.ResponseWriter, r *http.Request) {
	var input map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Cant create project", http.StatusBadRequest)
		return
	}
	
	project :=  &models.Project{
		Name: input["name"].(string),
		Description: input["description"].(string),
		Public: input["public"].(bool),
	}

	if err := database.DB.Create(project).Error; err != nil {
		http.Error(w, "Couldn't create project", http.StatusInternalServerError)
		return
	}

	if err := os.MkdirAll("./uploads/" + strconv.Itoa(int(project.Id)), 0755); err != nil {
		http.Error(w, "Couldn't create project directory", http.StatusInternalServerError)
		return
	}

	type NewProject struct {
		Id 			uint   `json:"id"`
		Name 		string `json:"name"`
		Description string `json:"description"`
		Public 		bool   `json:"public"`
	}

	newProject := NewProject{
		Id:          project.Id,
		Name:        project.Name,
		Description: project.Description,
		Public:      project.Public,
	}

	access := &models.Access{
		UserId:    uint(input["user_id"].(float64)),
		ProjectId: project.Id,
	}

	if err := database.DB.Create(access).Error; err != nil {
		http.Error(w, "Couldn't create access on project", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"message": "Created project",
		"data": newProject,
	})	
}

