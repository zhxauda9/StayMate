package service

import (
	"errors"

	"github.com/zhxauda9/StayMate/internal/dal/postgres"
	"github.com/zhxauda9/StayMate/models"
)

type UserService interface {
	CreateUser(user models.User) error
	GetUserByID(id int) (models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(id int, user models.User) error
	DeleteUser(id int) error
}

type userService struct {
	repo postgres.UserRepo
}

func NewUserService(repo postgres.UserRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(user models.User) error {
	if user.Name == "" || user.Email == "" {
		return errors.New("name and email cannot be empty")
	}
	return s.repo.CreateUser(user)
}

func (s *userService) GetUserByID(id int) (models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *userService) UpdateUser(id int, user models.User) error {
	if user.Name == "" && user.Email == "" {
		return errors.New("nothing to update, name or email must be provided")
	}
	return s.repo.UpdateUser(id, user)
}

func (s *userService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
