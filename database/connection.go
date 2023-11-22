package database

import (
	"fmt"
	"github.com/kapalfa/go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
);

var DB *gorm.DB

func ConnectDB() {
	// comments
	var err error 
	dsn := "host=localhost port=5432 user=postgres password=pass dbname=mydb sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Could not connect to the database")
		return
	}

	fmt.Println("Database connection successfully opened")

	err = DB.AutoMigrate(&models.User{})	
	if err != nil {
		fmt.Println("Could not migrate the database")
	}

}