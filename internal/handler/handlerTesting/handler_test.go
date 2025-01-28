package handlertesting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zhxauda9/StayMate/internal/dal/postgres"
	"github.com/zhxauda9/StayMate/internal/handler"
	"github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/models"
	"gorm.io/gorm"
)

// Mock AuthService
type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockAuthService) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *mockAuthService) ValidateToken(token string) (bool, error) {
	args := m.Called(token)
	return args.Bool(0), args.Error(1)
}

func (m *mockAuthService) GetUserFromToken(token string) (models.User, error) {
	args := m.Called(token)
	return args.Get(0).(models.User), args.Error(1)
}

// Test for Register method
func TestRegister(t *testing.T) {
	myLogger.Log = myLogger.NewZeroLoggerV3()
	mockService := new(mockAuthService)
	db := new(gorm.DB)
	verifyRepo := postgres.NewVerifyRepository(db)
	handler := handler.NewAuthHandler(mockService, verifyRepo)

	// Test case 1: Valid user
	user := models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "Aa123456$",
	}
	mockService.On("Register", user).Return(nil)

	t.Log("Running Test: Valid User Registration")
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]string
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "User created successfully", response["message"])

	// Test case 2: Invalid user (missing email)
	invalidUser := models.User{
		Name: "Invalid User",
	}
	mockService.On("Register", invalidUser).Return(fmt.Errorf("invalid input"))

	t.Log("Running Test: Invalid User Registration (missing email)")
	body, _ = json.Marshal(invalidUser)
	req = httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	rec = httptest.NewRecorder()

	handler.Register(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}

// Test for Login method
func TestLogin(t *testing.T) {
	myLogger.Log = myLogger.NewZeroLoggerV3()
	mockService := new(mockAuthService)
	db := new(gorm.DB)
	verifyRepo := postgres.NewVerifyRepository(db)
	handler := handler.NewAuthHandler(mockService, verifyRepo)

	credentials := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    "test@example.com",
		Password: "Aa123456$",
	}

	token := "mocked-token"
	mockService.On("Login", credentials.Email, credentials.Password).Return(token, nil)

	t.Log("Running Test: Valid Login")
	body, _ := json.Marshal(credentials)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]string
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Login successful", response["message"])
	assert.Equal(t, token, response["token"])

	// Test case 2: Invalid Login (wrong password)
	mockService.On("Login", credentials.Email, "wrong-password").Return("", fmt.Errorf("invalid credentials"))

	t.Log("Running Test: Invalid Login (wrong password)")
	credentials.Password = "wrong-password"
	body, _ = json.Marshal(credentials)
	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	rec = httptest.NewRecorder()

	handler.Login(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	mockService.AssertExpectations(t)
}

// Test for ValidateToken method
func TestValidateToken(t *testing.T) {
	myLogger.Log = myLogger.NewZeroLoggerV3()
	mockService := new(mockAuthService)
	db := new(gorm.DB)
	verifyRepo := postgres.NewVerifyRepository(db)
	handler := handler.NewAuthHandler(mockService, verifyRepo)

	token := "mocked-token"
	mockService.On("ValidateToken", token).Return(true, nil)

	t.Log("Running Test: Valid Token")
	req := httptest.NewRequest(http.MethodGet, "/auth/validate", nil)
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: token})
	rec := httptest.NewRecorder()

	handler.ValidateToken(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]string
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "I'm logging!", response["message"])
	assert.Equal(t, token, response["token"])

	// Test case 2: Invalid Token
	mockService.On("ValidateToken", "invalid-token").Return(false, fmt.Errorf("invalid token"))

	t.Log("Running Test: Invalid Token")
	req = httptest.NewRequest(http.MethodGet, "/auth/validate", nil)
	req.AddCookie(&http.Cookie{Name: "Authorization", Value: "invalid-token"})
	rec = httptest.NewRecorder()

	handler.ValidateToken(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	mockService.AssertExpectations(t)
}
