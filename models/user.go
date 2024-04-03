package models

import "time"

type User struct {
	Id                uint      `gorm:"primaryKey" json:"id"`
	Email             string    `gorm:"unique" json:"email"`
	Password          string    `json:"password"`
	RefreshToken      string    `json:"refresh_token"`
	GithubToken       string    `json:"github_token"`
	Verified          bool      `json:"verified" gorm:"default:false"`
	VerificationToken string    `json:"verification_token"`
	ResetToken        string    `json:"reset_token"`
	ResetTokenExpires time.Time `json:"reset_token_expires"`
	Accesses          []Access  `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"accesses"`
}
