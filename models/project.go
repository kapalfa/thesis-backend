package models

type Project struct {
	Id 			 uint   `gorm:"primaryKey" json:"id"`
	Name	 	 string `gorm:"not null" json:"name"`
	Description  string `gorm:"not null" json:"description"`
	Public 		 bool   `json:"public"`
	Accesses	[]Access `gorm:"foreignKey:ProjectId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"accesses"`
}
