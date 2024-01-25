package projectsCRUD 

import (
	"encoding/json"
	"net/http"
	"github.com/kapalfa/go/models"
	"github.com/kapalfa/go/database"
)

type Collaborator struct {
	UserId uint `json:"user_id"`
	Email string `json:"email"`
}

type PublicProjectResponse struct {
	Id uint `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Public bool `json:"public"`
	Collaborators []Collaborator `json:"collaborators"`
}

func GetPublicProjects(w http.ResponseWriter, r *http.Request) {
	var projects []models.Project
	var responses []PublicProjectResponse
	database.DB.Model(&models.Project{}).Where("public = ?", true).Find(&projects)
	for _, project := range projects {
		var collaborators []Collaborator
		database.DB.Model(&models.Access{}).Select("user_id, email").Joins("left join users on accesses.user_id = users.id").Where("project_id = ?", project.Id).Scan(&collaborators)
		
		response := PublicProjectResponse{
			Id: project.Id,
			Name: project.Name,
			Description: project.Description,
			Public: project.Public,
			Collaborators: collaborators,
		}
		responses = append(responses, response)	
	}
	json.NewEncoder(w).Encode(responses)
}