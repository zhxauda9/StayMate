package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/models"
	"gorm.io/gorm"
)

type verifyHandler struct {
	db *gorm.DB
}

func NewVerifyHandler(db *gorm.DB) *verifyHandler {
	return &verifyHandler{db: db}
}

// Verify handles the email verification process
func (h *verifyHandler) Verify(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming JSON body
	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	// Decode the request body into the struct
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		myLogger.Log.Error().Err(err).Msg("Error decoding JSON request body")
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Check if email and code are provided
	if request.Email == "" || request.Code == "" {
		http.Error(w, "Email and verification code are required", http.StatusBadRequest)
		return
	}

	// Check the UsersEmailConfirm table for the provided email and code
	var userEmailConfirm models.UsersEmailConfirm
	err = h.db.Where("email = ? AND code = ?", request.Email, request.Code).First(&userEmailConfirm).Error
	if err != nil {
		myLogger.Log.Error().Err(err).Str("email", request.Email).Msg("Failed to find matching email/code in UsersEmailConfirm")
		http.Error(w, "Invalid verification code or email", http.StatusBadRequest)
		return
	}

	// If email and code match, update the status of the user to "active"
	var user models.User
	err = h.db.Where("email = ?", request.Email).First(&user).Error
	if err != nil {
		myLogger.Log.Error().Err(err).Str("email", request.Email).Msg("Failed to find user in Users table")
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}

	// Update user status to "active"
	user.Status = "active"
	err = h.db.Save(&user).Error
	if err != nil {
		myLogger.Log.Error().Err(err).Str("email", request.Email).Msg("Failed to update user status")
		http.Error(w, "Failed to update user status", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email successfully verified, user status is now active"))
}
