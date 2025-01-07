package postgres

import (
	"errors"
	"fmt"

	"github.com/zhxauda9/StayMate/models"
	"gorm.io/gorm"
)

type UserRepo interface {
	CreateUser(user models.User) error
	GetUserByID(id int) (models.User, error)
	GetAllUsers(sort, filterStart, filterEnd string, limit, offset int) ([]models.User, error)
	UpdateUser(id int, user models.User) error
	DeleteUser(id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepo {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user models.User) error {
	if user.Name == "" || user.Email == "" || user.Status == "" {
		return errors.New("name and email fields cannot be empty")
	}
	if err := r.db.Create(&user).Error; err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

func (r *userRepository) GetUserByID(id int) (models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.User{}, fmt.Errorf("user with ID %d not found", id)
		}
		return models.User{}, fmt.Errorf("error fetching user by ID: %v", err)
	}
	return user, nil
}

func (r *userRepository) GetAllUsers(sort, filterStart, filterEnd string, limit, offset int) ([]models.User, error) {
	var users []models.User
	query := r.db
	if sort != "" {
		query = query.Order(sort)
	}
	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("error fetching all users: %v", err)
	}

	query1 := r.db.Model(&models.User{})
	if filterStart != "" && filterEnd != "" {
		query1 = query1.Where("price >= ? AND check_in <= ?", filterStart, filterEnd)
	}
	if err := query1.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("error fetching filtered users: %v", err)
	}

	return users, nil
}

func (r *userRepository) UpdateUser(id int, user models.User) error {
	var existingUser models.User
	if err := r.db.First(&existingUser, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user with ID %d not found", id)
		}
		return fmt.Errorf("error fetching user for update: %v", err)
	}
	if user.Name != "" {
		existingUser.Name = user.Name
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.Status != "" {
		existingUser.Status = user.Status
	}
	if err := r.db.Save(&existingUser).Error; err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}
	return nil
}

func (r *userRepository) DeleteUser(id int) error {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user with ID %d not found", id)
		}
		return fmt.Errorf("error fetching user for deletion: %v", err)
	}
	if err := r.db.Delete(&user).Error; err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}
