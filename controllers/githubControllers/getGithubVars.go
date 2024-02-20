package githubControllers

import "os"

func GetGithubClientID() string {
	githubClientID := os.Getenv("CLIENT_ID")
	return githubClientID
}

func GetGithubClientSecret() string {
	githubClientSecret := os.Getenv("CLIENT_SECRET")
	return githubClientSecret
}
