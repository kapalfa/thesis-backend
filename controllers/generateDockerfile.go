package controllers 

import (
	"fmt"
	"os"
)

func GenerateDockerfile(language string, filename string, command string) {
	var baseImage string 

	switch language {
	case "c", "cpp":
		baseImage = "gcc:latest"
	case "java":
		baseImage = "openjdk:latest"
	case "python":
		baseImage = "python:latest"
	default:
		fmt.Printf("Unsupported language: %s\n", language)
		return 
	}

	dockerfileTemplate:= `
	FROM  %s
	WORKDIR /app
	COPY ../%s /app
	RUN 
	CMD ["%s"]
	`

	dockerfile := fmt.Sprintf(dockerfileTemplate, baseImage, filename, command)

	if err := os.WriteFile("./DockerContents/Dockerfile", []byte(dockerfile), 0644); err != nil {
		fmt.Printf("Error writing Dockerfile: %v\n", err)
		return 
	}


}