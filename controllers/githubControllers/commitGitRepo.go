package githubControllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"golang.org/x/oauth2"
)

func CommitGitRepo(w http.ResponseWriter, r *http.Request) {
	type CommitGitRepoRequest struct {
		ProjectId string `json:"projectid"`
		UserId    uint   `json:"userid"`
	}
	var request CommitGitRepoRequest
	err := json.NewDecoder(r.Body).Decode(&request) // get project and user id
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var user models.User
	database.DB.First(&user, request.UserId)
	accessToken := user.GithubToken // get github token from user model

	var proj models.Project
	projid, err := strconv.ParseUint(request.ProjectId, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	database.DB.First(&proj, projid)
	repo := proj.Name // get project name from project model

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+string(accessToken))
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
	branch := "main"                                 // default branch
	projectFolder := string(request.ProjectId) + "/" // project folder path
	command := "https://api.github.com/repos/" + username + "/" + repo + "/commits"

	req, err = http.NewRequest("GET", command, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Printf("Error making request to fetch commits: %v", err)
		return
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}
	type Commit struct {
		Sha string `json:"sha"`
	}
	var commits []Commit
	err = json.Unmarshal(b, &commits)
	if err != nil {
		log.Printf("Error parsing the response body: %v", err)
		return
	}
	var sha string
	if len(commits) > 0 {
		sha = commits[0].Sha
	} else {
		sha = ""
	}

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	ghClient := github.NewClient(oauthClient)
	tree, err := CreateGitTree(ghClient, accessToken, username, repo, branch, projectFolder)
	if err != nil {
		log.Printf("Error creating git tree: %v", err)
		return
	}

	commit, _, err := ghClient.Git.CreateCommit(context.Background(), username, repo, &github.Commit{
		Message: github.String("Initial commit"),
		Tree:    tree,
		Parents: []github.Commit{{SHA: &sha}}, // Fix: Pass the address of the sha variable
	})
	if err != nil {
		log.Printf("Error creating commit: %v", err)
		return
	}

	ref, _, err := ghClient.Git.UpdateRef(context.Background(), username, repo, &github.Reference{
		Ref: github.String("refs/heads/" + branch),
		Object: &github.GitObject{
			SHA: commit.SHA,
		},
	}, false)
	if err != nil {
		log.Printf("Error updating ref: %v", err)
		return
	}
	log.Println("Ref updated: ", ref.Object.SHA)
	// 	cmd := exec.Command("bash", "./script/createTree.sh",
	// 		"-t", accessToken,
	// 		"-u", username,
	// 		"-r", repo,
	// 		"-b", branch,
	// 		"-p", projectFolder,
	// 		"-s", sha,
	// 	)
	// 	output, err := cmd.CombinedOutput()
	// 	if err != nil {
	// 		fmt.Printf("Error running the script: %s\n", err)
	// 		fmt.Println("Output: ", string(output))
	// 		return
	// 	}
	// 	fmt.Println("Output: ", string(output))
}
