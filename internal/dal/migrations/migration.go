package migrations

import (
	"github.com/zhxauda9/StayMate/models"
	"gorm.io/gorm"
)

func AutoMigrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(&models.Booking{}, &models.Room{Number: 14, Class: "lux", Price: 4999}, &models.User{})
}
