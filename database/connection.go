package database

import (
	"fmt"
	"os"

	"github.com/kapalfa/go/migrations"
	"github.com/kapalfa/go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	//"github.com/joho/godotenv"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	//	err = godotenv.Load()
	//	if err != nil {
	//		fmt.Println("Error loading .env file")
	//	}
	dbport := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	host := os.Getenv("DB_HOST")
	sslrootcert := os.Getenv("DB_SSLROOTCERT")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s", host, dbport, user, password, dbname, sslmode, sslrootcert)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Could not connect to the database")
		return
	}
	fmt.Println("Database connection successfully opened")

	migration := migrations.Migration1617448756{}
	err = migration.Migrate(DB)
	if err != nil {
		fmt.Println("Could not migrate the database")
		return
	}
	err = DB.AutoMigrate(&models.User{}, &models.Project{}, &models.Access{})
	if err != nil {
		fmt.Println("Could not migrate the database")
	}

}
