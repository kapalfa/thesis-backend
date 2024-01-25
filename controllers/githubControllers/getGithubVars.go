package githubControllers

import (
	"log"
	"github.com/joho/godotenv"
	"os"
)

func GetGithubClientID() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	githubClientID := os.Getenv("CLIENT_ID")
	return githubClientID
}

func GetGithubClientSecret() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	githubClientSecret := os.Getenv("CLIENT_SECRET")
	return githubClientSecret
}