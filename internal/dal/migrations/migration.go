package migrations

import (
	"github.com/zhxauda9/StayMate/models"
	"gorm.io/gorm"
)

func AutoMigrateDatabase(db *gorm.DB) error {
	db.AutoMigrate(&models.Booking{}, &models.Room{}, &models.User{})
	Fill(db)
	return nil
}
