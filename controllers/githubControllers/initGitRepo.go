package githubControllers	

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)
func InitGitRepo(w http.ResponseWriter, r *http.Request) {
	type InitGitRepoRequest struct {
		UserId 		uint 	`json:"userid"`
		ProjectId 	uint 	`json:"projectid"`
	}

	var request InitGitRepoRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	database.DB.First(&user, request.UserId)
	if user.Id == 0 {
		log.Println("User not found")
		return
	}

	githubToken := user.GithubToken	
	if githubToken == "" {
		log.Println("Github token not found")
		return
	}

	var proj models.Project
	database.DB.First(&proj, request.ProjectId)
	if proj.Id == 0 {
		log.Println("Project not found")
		return
	}	

	s := strings.NewReader(`{"name": "` + proj.Name + `", "description": "` + proj.Description + `", "private": true, "auto_init": true}`)
	req, err := http.NewRequest("POST", "https://api.github.com/user/repos", s)
	if err != nil {
		log.Printf("Error creating request to initialize repo: %v", err)
		return
	}

	req.Header.Set("Authorization", "Bearer " + string(githubToken))	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making request to initialize repo: %v", err)
		return
	}
	defer resp.Body.Close()
}