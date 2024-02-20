package githubControllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"github.com/kapalfa/go/config"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

func getContents(path, username, repoName, accessToken string, projectId uint, ctx context.Context, bkt *storage.BucketHandle) {
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+username+"/"+repoName+"/contents"+path, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
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
			dirName := strconv.Itoa(int(projectId)) + "/" + itemPath + "/"
			log.Println("Creating directory: ", dirName)
			obj := bkt.Object(dirName)
			wc := obj.NewWriter(ctx)
			if err := wc.Close(); err != nil {
				log.Printf("Error creating directory: %v", err)
				return
			}
			getContents(itemPath, username, repoName, accessToken, projectId, ctx, bkt)
		} else {
			req, _ = http.NewRequest("GET", "https://api.github.com/repos/"+username+"/"+repoName+"/contents"+itemPath, nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)
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
			filename := strconv.Itoa(int(projectId)) + "/" + itemPath
			obj := bkt.Object(filename)
			wc := obj.NewWriter(ctx)
			if _, err := wc.Write(decoded); err != nil {
				log.Printf("Error writing file: %v", err)
				return
			}
			if err := wc.Close(); err != nil {
				log.Printf("Error closing writer: %v", err)
				return
			}
		}
	}
}

type Req struct {
	RepoName string `json:"repoName"`
}

func DownloadGitRepo(w http.ResponseWriter, r *http.Request) {
	ctx := config.Ctx
	bkt := config.Bucket

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
		Name:        myreq.RepoName,
		Description: "test", // get from user
		Public:      false,
	}
	if err := database.DB.Create(project).Error; err != nil {
		http.Error(w, "Couldn't create project", http.StatusInternalServerError)
		return
	}
	access := &models.Access{ // create new access
		UserId:    user.Id,
		ProjectId: project.Id,
	}
	if err := database.DB.Create(access).Error; err != nil {
		http.Error(w, "Couldn't create access on project", http.StatusInternalServerError)
		return
	}
	dirName := strconv.Itoa(int(project.Id)) + "/"
	obj := bkt.Object(dirName)
	wc := obj.NewWriter(ctx)
	if err := wc.Close(); err != nil {
		http.Error(w, "Couldn't create project directory", http.StatusInternalServerError)
		return
	}

	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil) // get username
	req.Header.Set("Authorization", "Bearer "+accessToken)
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

	getContents("", username, myreq.RepoName, accessToken, project.Id, ctx, bkt) // get contents of repo
}
