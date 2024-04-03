package models

type Access struct {
	UserId 		uint 	`gorm:"primaryKey" json:"user_id"`
	ProjectId 	uint 	`gorm:"primaryKey" json:"project_id"`
	Status 		string 	`gorm:"type:access_type;default:'accepted'" json:"status"`
	User 		User 	`gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	Project 	Project `gorm:"foreignKey:ProjectId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"project"`
}