package servicetesting

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/zhxauda9/StayMate/internal/service"
	"github.com/zhxauda9/StayMate/models"
)

func TestValidateToken(t *testing.T) {
	// Настройка секретного ключа
	secret := "test_secret"
	os.Setenv("SECRET", secret)

	// Генерация токенов для тестов
	validToken, _ := GenerateToken(secret, jwt.MapClaims{"user_id": 1, "email": "test@example.com"}, time.Minute*10)
	expiredToken, _ := GenerateToken(secret, jwt.MapClaims{"user_id": 1, "email": "test@example.com"}, -time.Minute*10)
	invalidSignatureToken, _ := GenerateToken("wrong_secret", jwt.MapClaims{"user_id": 1, "email": "test@example.com"}, time.Minute*10)
	tokenWithNoExp, _ := GenerateToken(secret, jwt.MapClaims{"user_id": 1, "email": "test@example.com"}, 0)

	// Подготовка входных данных
	inputs := []struct {
		name     string
		token    string
		expected bool
	}{
		// Корректные сценарии
		{"Valid token", validToken, true},

		// Ошибки формата и подписи
		{"Empty token", "", false},
		{"Malformed token", "invalid.token.structure", false},
		{"Token with invalid signature", invalidSignatureToken, false},

		// Ошибки срока действия
		{"Expired token", expiredToken, false},
		{"Token without expiration", tokenWithNoExp, false},

		// Проверка алгоритма подписи
		{"Token with invalid signing method", GenerateTokenWithAlgorithm(secret, jwt.MapClaims{"user_id": 1}, jwt.SigningMethodRS256, time.Minute*10), false},
	}

	// Подготовка фейкового сервиса пользователя
	fakeUserServ := &fakeUserService{}
	authService := service.NewAuthService(fakeUserServ)

	// Запуск тестов
	for _, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			result, err := authService.ValidateToken(input.token)
			t.Logf("Calling ValidateToken(%v), Expected: %v, Actual: %v", input.token, input.expected, result)
			if result != input.expected {
				t.Errorf("Test %v failed. Expected: %v, Got: %v, Error: %v", input.name, input.expected, result, err)
			}
		})
	}
}

// GenerateToken creates a JWT token with the provided secret, claims, and expiration duration
func GenerateToken(secret string, claims jwt.MapClaims, duration time.Duration) (string, error) {
	if duration > 0 {
		claims["exp"] = time.Now().Add(duration).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateTokenWithAlgorithm(secret string, claims jwt.MapClaims, method jwt.SigningMethod, duration time.Duration) string {
	if duration > 0 {
		claims["exp"] = time.Now().Add(duration).Unix()
	}
	token := jwt.NewWithClaims(method, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

type fakeUserService struct{}

func (f *fakeUserService) CreateUser(user models.User) error {
	return nil
}

func (f *fakeUserService) GetUserByID(id int) (models.User, error) {
	return models.User{}, nil
}

func (f *fakeUserService) GetUserByEmail(email string) (models.User, error) {
	if email == "test@example.com" {
		return models.User{ID: 1, Email: email, Password: "$2a$10$XyqW8XxQwz/vPC8PbO/0CuODDZUkv1mjcOoR3r.n/VO9XfyO8W19G"}, nil // bcrypt hash for "password"
	}
	return models.User{}, errors.New("user not found")
}

func (f *fakeUserService) GetAllUsers(sort string, page int) ([]models.User, error) {
	return nil, nil
}

func (f *fakeUserService) UpdateUser(id int, user models.User) error {
	return nil
}

func (f *fakeUserService) DeleteUser(id int) error {
	return nil
}
