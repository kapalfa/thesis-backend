package invitationsCRUD 
import (
	"encoding/json"
	"net/http"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"github.com/kapalfa/go/functions"
)
type CollaboratorsEmail struct {
	ProjectId uint `json:"project_id"`
	Email string `json:"email"`
}
type ProjectInfo struct {
	ProjectId 			uint `json:"project_id"`
	ProjectName 		string `json:"project_name"`
	ProjectDescription	string `json:"project_description"`
	CollaboratorsEmail []string `json:"collaborators"`
}
func GetInvitations(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	id, err := functions.GetUserId(auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var projectIds []uint
	var projects []models.Project
	var projectInfos []ProjectInfo
	var collaboratorsEmail []CollaboratorsEmail
	database.DB.Model(&models.Access{}).Where("user_id = ? AND status = ?", id, "pending").Pluck("project_id", &projectIds)
	if len(projectIds) == 0 {
		json.NewEncoder(w).Encode(projectInfos)
		return
	}
	database.DB.Table("accesses").Select("accesses.project_id, users.email").Joins("left join users ON accesses.user_id = users.id").Where("accesses.project_id IN ? AND accesses.status = ?", projectIds, "accepted").Scan(&collaboratorsEmail)
	database.DB.Model(&models.Project{}).Where("id IN ?", projectIds).Find(&projects)
	for _, project := range projects {
		var email []string 
		for _, collaborator := range collaboratorsEmail {
			if collaborator.ProjectId == project.Id {
				email = append(email, collaborator.Email)
			}
		}
		projectInfos = append(projectInfos, ProjectInfo{ProjectId: project.Id, ProjectName: project.Name, ProjectDescription: project.Description, CollaboratorsEmail: email})
	}
	json.NewEncoder(w).Encode(projectInfos)
}
