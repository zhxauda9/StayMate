package handlertesting

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zhxauda9/StayMate/internal/handler"
	"github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/models"
)

// Mock UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetUserByID(id int) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(email string) (models.User, error) {
	args := m.Called(email)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserService) GetAllUsers(sort string, page int) ([]models.User, error) {
	args := m.Called(sort, page)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(id int, user models.User) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test for CreateUser
func TestCreateUser(t *testing.T) {
	myLogger.Log = myLogger.NewZeroLogger()

	mockService := new(MockUserService)
	handler := handler.NewUserHandler(mockService)

	newUser := models.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepassword",
		Status:   "active",
	}

	mockService.On("CreateUser", newUser).Return(nil)

	body, _ := json.Marshal(newUser)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.CreateUser(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var response models.User
	_ = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, newUser.Name, response.Name)
	mockService.AssertExpectations(t)
}

// Test for GetAllUsers
func TestGetAllUsers(t *testing.T) {
	mockService := new(MockUserService)
	handler := handler.NewUserHandler(mockService)

	users := []models.User{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
	}

	mockService.On("GetAllUsers", "", 1).Return(users, nil)

	req := httptest.NewRequest(http.MethodGet, "/users?page=1", nil)
	rec := httptest.NewRecorder()

	handler.GetUsers(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var response []models.User
	_ = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Len(t, response, 2)
	assert.Equal(t, users[0].Name, response[0].Name)
	mockService.AssertExpectations(t)
}

// // Test for GetUserByID
// func TestGetUserByID(t *testing.T) {
// 	mockService := new(MockUserService)
// 	handler := handler.NewUserHandler(mockService)

// 	testUser := models.User{
// 		ID:    1,
// 		Name:  "Jane Doe",
// 		Email: "jane@example.com",
// 	}

// 	mockService.On("GetUserByID", testUser.ID).Return(testUser, nil)

// 	req := httptest.NewRequest(http.MethodGet, "/users/{1}", nil)
// 	rec := httptest.NewRecorder()

// 	handler.GetUserByID(rec, req)

// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	var response models.User
// 	_ = json.Unmarshal(rec.Body.Bytes(), &response)
// 	assert.Equal(t, testUser.Name, response.Name)
// 	mockService.AssertExpectations(t)
// }

// Test for GetUserByID
// func TestUpdateUser(t *testing.T) {
// 	// Initialize the handler with a mock service
// 	mockService := new(MockUserService)
// 	handler := handler.NewUserHandler(mockService)

// 	// User data to be updated
// 	updatedUser := models.User{
// 		Name:  "Updated User",
// 		Email: "updated@example.com",
// 	}

// 	// Marshal the updated user data to JSON
// 	body, _ := json.Marshal(updatedUser)

// 	// Create a request with the method PUT and pass the updated user data
// 	req := httptest.NewRequest(http.MethodPut, "/users/123", bytes.NewReader(body))
// 	rec := httptest.NewRecorder()

// 	// Call the handler's UpdateUser method
// 	handler.UpdateUser(rec, req)

// 	// Assert the response status and body
// 	assert.Equal(t, http.StatusOK, rec.Code)

// 	var updatedUserResponse models.User
// 	err := json.NewDecoder(rec.Body).Decode(&updatedUserResponse)
// 	assert.NoError(t, err)
// 	assert.Equal(t, updatedUser.Name, updatedUserResponse.Name)
// 	assert.Equal(t, updatedUser.Email, updatedUserResponse.Email)
// }

// // Test for DeleteUser
// func TestDeleteUser(t *testing.T) {
// 	mockService := new(MockUserService)
// 	handler := handler.NewUserHandler(mockService)

// 	mockService.On("DeleteUser", 1).Return(nil)

// 	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
// 	rec := httptest.NewRecorder()

// 	handler.DeleteUser(rec, req)

// 	assert.Equal(t, http.StatusNoContent, rec.Code)
// 	mockService.AssertExpectations(t)
// }
