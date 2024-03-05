package config

import (
	"os"
)

func Config(key string) string {
	//err := godotenv.Load(".env")
	//if err != nil {
	//	fmt.Println("Error loading .env file")
	//}
	return os.Getenv(key)
}
