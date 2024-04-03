package githubControllers

import (
	"fmt"
	"net/http"
)

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := GetGithubClientID()
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&scope=repo%%20user",
		githubClientID,
		"https://thesis-frontend-snowy.vercel.app/github/callback",
		"randomString",
	)
	w.Write([]byte(redirectURL))
}
