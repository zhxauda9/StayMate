package models

import (
	"time"

	"gorm.io/gorm"
)

type UsersEmailConfirm struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey"`     // Primary Key
	Email     string    `gorm:"uniqueIndex"`    // Unique Email
	Code      string    `gorm:"not null"`       // Confirmation Code
	CreatedAt time.Time `gorm:"autoCreateTime"` // Created At with auto timestamp
}
