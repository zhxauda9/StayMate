package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"not null;unique"`
	Password string `json:"password"`
	Role     string `json:"role" gorm:"default:user"`
	Status   string `json:"status" gorm:"not null;default:inactive"`
	Photo    string `json:"photo" gorm:"default:static/pictures/default/user.jpg"`

	Chats []Chat `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
