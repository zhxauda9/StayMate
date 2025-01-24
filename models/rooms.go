package models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Number      int     `json:"number" gorm:"not null"`
	Class       string  `json:"class" gorm:"not null"`
	Price       float64 `json:"price" gorm:"not null"`
	Status      string  `json:"status" gorm:"not null"`
	Photo       string  `json:"photo" gorm:"default:static/pictures/default/hotel.jpg"`
	Description string  `json:"description" gorm:"default:\"VERY BEAUITFUUUUULLLLLLLL\""`
}
