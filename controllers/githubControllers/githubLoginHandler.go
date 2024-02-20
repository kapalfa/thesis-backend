package githubControllers

import (
	"fmt"
	"net/http"
	"log"
)

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("IN GITHUB LOGIN HANDLER")
	githubClientID := GetGithubClientID()
	redirectURL := fmt.Sprintf(
	 	"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&scope=repo%%20user",
	 	githubClientID,
	 	"https://localhost:5173/github/callback",
		"randomString",
	)
    w.Write([]byte(redirectURL))
}