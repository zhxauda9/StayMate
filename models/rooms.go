package models

import _ "gorm.io/gorm"

type Room struct {
	ID     int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Number int     `json:"number" gorm:"not null"`
	Class  string  `json:"type" gorm:"not null"`
	Price  float64 `json:"price" gorm:"not null"`
}
