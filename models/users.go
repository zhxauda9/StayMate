package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name" gorm:"not null"`
	Email  string `json:"email" gorm:"not null"`
	Status string `json:"status" gorm:"not null"`
}
