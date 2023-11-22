package models

type User struct {
	Id 			 uint   `gorm:"primaryKey" json:"id"`
	Name	 	 string `gorm:"unique" json:"name"`
	Email    	 string `gorm:"unique" json:"email"`
	Password 	 string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}