package models

type User struct {
	Id 			 uint   `gorm:"primaryKey" json:"id"`
	Email    	 string `gorm:"unique" json:"email"`
	Password 	 string `json:"password"`
	RefreshToken string `json:"refresh_token"`
	Accesses 	 []Access `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"accesses"`
}