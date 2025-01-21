package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/zhxauda9/StayMate/models"
)

type AuthService interface {
	Register(user models.User) error
	Login(email, password string) (string, error)
	ValidateToken(token string) (bool, error)
	GetUserFromToken(tokenString string) (models.User, error)
}

type authService struct {
	userService UserService // Используем сервис пользователей
}

func NewAuthService(userService UserService) AuthService {
	return &authService{userService: userService}
}

func (s *authService) Register(user models.User) error {
	return s.userService.CreateUser(user)
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userService.GetUserByEmail(email)
	if err != nil {
		return "Error: invalid email or password", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("Invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, errors.New("token is empty")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return false, errors.New("missing expiration claim")
		}
		if time.Now().Unix() > int64(exp) {
			return false, errors.New("token expired")
		}

		return true, nil
	}

	return false, errors.New("invalid token")
}

func (s *authService) GetUserFromToken(tokenString string) (models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil || !token.Valid {
		return models.User{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return models.User{}, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return models.User{}, errors.New("missing user_id in token")
	}

	user, err := s.userService.GetUserByID(int(userID))
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
