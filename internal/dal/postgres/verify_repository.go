package postgres

import (
	"github.com/zhxauda9/StayMate/models"
	"gorm.io/gorm"
)

type VerifyRepository struct {
	db *gorm.DB
}

func NewVerifyRepository(db *gorm.DB) *VerifyRepository {
	return &VerifyRepository{db: db}
}

func (r *VerifyRepository) InsertCode(row models.UsersEmailConfirm) error {
	// Insert the row into the UsersEmailConfirm table
	if err := r.db.Create(&row).Error; err != nil {
		return err
	}
	return nil
}
