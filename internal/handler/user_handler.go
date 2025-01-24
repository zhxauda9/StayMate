package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-playground/validator/v10"
	l "github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/internal/service"
	"github.com/zhxauda9/StayMate/models"
)

type userHandler struct {
	userService service.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to create a new user")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	file, fileheader, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Error uploading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	photoPath := "/static/pictures/storage" + filepath.Base(fileheader.Filename)
	outFile, err := os.Create(photoPath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		l.Log.Error().Err(err).Msg("Error decoding user data")
		http.Error(w, fmt.Sprintf("Error decoding user data: %v", err), http.StatusBadRequest)
		return
	}

	user.Photo = photoPath

	err = h.userService.CreateUser(user)
	if err != nil {
		l.Log.Error().Msg(fmt.Sprintf("Error creating user: %v", err))
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}
	l.Log.Info().Msg("User created successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *userHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to fetch all users")

	sort := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	users, err := h.userService.GetAllUsers(sort, page)
	if err != nil {
		l.Log.Error().Msg(fmt.Sprintf("Error fetching users: %v", err))
		http.Error(w, fmt.Sprintf("Error fetching users: %v", err), http.StatusInternalServerError)
		return
	}
	l.Log.Debug().Msg("Fetched all users successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *userHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to fetch a specific user")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		l.Log.Warn().Msg(fmt.Sprintf("Invalid user ID: %v", idStr))
		http.Error(w, fmt.Sprintf("Error parsing user ID: %v", err), http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		l.Log.Error().Msg(fmt.Sprintf("Error fetching user with ID %d: %v", id, err))
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusNotFound)
		return
	}
	l.Log.Info().Int("UserID", id).Msg("Fetched user successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to update a user")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		l.Log.Warn().Msg(fmt.Sprintf("Invalid user ID: %v", idStr))
		http.Error(w, fmt.Sprintf("Error parsing user ID: %v", err), http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		l.Log.Error().Msg(fmt.Sprintf("Error decoding user data: %v", err))
		http.Error(w, fmt.Sprintf("Error decoding user data: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(user); err != nil {
		l.Log.Warn().Msg(fmt.Sprintf("Validation error: %v", err))
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}
	l.Log.Info().Msg("Validation successful for user update")

	err = h.userService.UpdateUser(id, user)
	if err != nil {
		l.Log.Error().Msg(fmt.Sprintf("Error updating user with ID %d: %v", id, err))
		http.Error(w, fmt.Sprintf("Error updating user: %v", err), http.StatusInternalServerError)
		return
	}
	l.Log.Info().Int("UserID", id).Msg("Updated user successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to delete a user")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		l.Log.Warn().Msg(fmt.Sprintf("Invalid user ID: %v", idStr))
		http.Error(w, fmt.Sprintf("Error parsing user ID: %v", err), http.StatusBadRequest)
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		l.Log.Error().Msg(fmt.Sprintf("Error deleting user with ID %d: %v", id, err))
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusInternalServerError)
		return
	}
	l.Log.Info().Int("UserID", id).Msg("Deleted user successfully")

	w.WriteHeader(http.StatusNoContent)
}
