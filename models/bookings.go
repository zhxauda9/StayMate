package models

import (
	"gorm.io/gorm"
	_ "gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID   int    `json:"user_id" gorm:"not null"`
	RoomID   int    `json:"room_id" gorm:"not null"`
	CheckIn  string `json:"check_in" gorm:"not null"`
	CheckOut string `json:"check_out" gorm:"not null"`
}
