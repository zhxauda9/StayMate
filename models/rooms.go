package models

import (
	"gorm.io/gorm"
	_ "gorm.io/gorm"
)

type Room struct {
	gorm.Model
	ID     int     `json:"id" gorm:"primaryKey"`
	Number int     `json:"number" gorm:"not null"`
	Class  string  `json:"type" gorm:"not null"`
	Price  float64 `json:"price" gorm:"not null"`
}
