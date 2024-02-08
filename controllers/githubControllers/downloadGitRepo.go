package githubControllers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"os"
	"io"
	"strconv"
)
func getContents(path, username, repoName, accessToken string, projectId uint) {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/" + username + "/" + repoName + "/contents" + path, nil)
	req.Header.Set("Authorization", "Bearer " + accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making request to fetch user email: %v", err)
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}
	var data []map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Printf("Error parsing the response body: %v", err)
		return
	}
	for _, item := range data {
		itemType, ok := item["type"].(string)
		if !ok {
			log.Printf("Cannot assert item[type] as string")
			continue
		}
		itemPath, ok := item["path"].(string)
		if !ok {
			log.Printf("Cannot assert item[path] as string")
			continue
		}
		if itemType == "dir" {
			err = os.MkdirAll("./uploads/" + strconv.Itoa(int(projectId)) + "/" + itemPath, 0755)
			getContents(itemPath, username, repoName, accessToken, projectId)
		} else {
			req, err = http.NewRequest("GET", "https://api.github.com/repos/" + username + "/" + repoName + "/contents" + itemPath, nil)
			req.Header.Set("Authorization", "Bearer " + accessToken)
			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("Error making request to fetch user email: %v", err)
				return
			}
			defer resp.Body.Close()			
			b, err = io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading response body: %v", err)
				return
			}
			var data2 map[string]interface{}
			err = json.Unmarshal(b, &data2)
			if err != nil {
				log.Printf("Error parsing the response body: %v", err)
				return
			}
			content, ok := data2["content"].(string)
			if !ok {
				log.Printf("Cannot assert data2[content] as string")
				continue
			}
			decoded, err := base64.StdEncoding.DecodeString(content)
			if err != nil {
				log.Printf("Error decoding content: %v", err)
				return
			}
			err = os.WriteFile("./uploads/" + strconv.Itoa(int(projectId)) + "/" + itemPath, decoded, 0644)
			if err != nil {
				log.Printf("Error writing file: %v", err)
				return
			}
		}
	}
}
type Req struct {
	RepoName string `json:"repoName"`
}
func DownloadGitRepo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userid"]
	var myreq Req
	err := json.NewDecoder(r.Body).Decode(&myreq)
	if err != nil {
		log.Println("Error decoding request: ", err)
		return
	}

	var user models.User
	database.DB.First(&user, userId) // find user by id
	accessToken := user.GithubToken

	project := &models.Project{ // create new project 
		Name: myreq.RepoName,
		Description: "test", // get from user 
		Public: false,
	}

	if err := database.DB.Create(project).Error; err != nil {
		http.Error(w, "Couldn't create project", http.StatusInternalServerError)
		return
	}

	access := &models.Access{ // create new access
		UserId: user.Id,
		ProjectId: project.Id,
	}

	if err := database.DB.Create(access).Error; err != nil {
		http.Error(w, "Couldn't create access on project", http.StatusInternalServerError)
		return
	}

	if err := os.MkdirAll("./uploads/" + strconv.Itoa(int(project.Id)), 0755); err != nil { // create project directory	
		http.Error(w, "Couldn't create project directory", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil) // get username 
	req.Header.Set("Authorization", "Bearer " + accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making request to fetch user email: %v", err)
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)

	if err != nil {
		log.Printf("Error parsing the response body: %v", err)
		return
	}
	username, ok := data["login"].(string)
	if !ok {
		log.Printf("Cannot assert data[login] as string")
		return
	}

	getContents("", username, myreq.RepoName, accessToken, project.Id) // get contents of repo

	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"status": "success",
	// 	"message": "Downloaded repo",
	// 	"data": newProject,
	// })
}