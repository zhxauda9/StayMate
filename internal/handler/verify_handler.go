package handler

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/internal/service"
	"github.com/zhxauda9/StayMate/models"
	"gorm.io/gorm"
)

type verifyHandler struct {
	db   *gorm.DB
	mail service.MailServiceImpl
}

func NewVerifyHandler(db *gorm.DB, mail service.MailServiceImpl) *verifyHandler {
	return &verifyHandler{db: db, mail: mail}
}

func (h *verifyHandler) SendVerifyCode(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		myLogger.Log.Error().Err(err).Msg("Error decoding JSON request body")
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	if request.Email == "" {
		myLogger.Log.Error().Msg("Email and verification code are required")
		http.Error(w, "Email and verification code are required", http.StatusBadRequest)
		return
	}

	code := GenerateCode()

	// Check if there's already a pending verification for the email
	var existingRecord models.UsersEmailConfirm
	err = h.db.Where("email = ?", request.Email).First(&existingRecord).Error
	if err == nil && time.Now().Before(existingRecord.CreatedAt.Add(5*time.Minute)) {
		myLogger.Log.Error().Msg("A verification code has already been sent. Please wait 5 minutes before requesting another.")
		http.Error(w, "A verification code has already been sent. Please wait 5 minutes before requesting another.", http.StatusTooManyRequests)
		return
	}

	// Create or update the record in the database
	newRecord := models.UsersEmailConfirm{
		Email:     request.Email,
		Code:      code,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute), // Set expiration time to 5 minutes from now
	}
	err = h.db.Where("email = ?", request.Email).Save(&newRecord).Error
	if err != nil {
		myLogger.Log.Error().Err(err).Msg("Error saving verification code to the database")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = h.mail.Send([]string{request.Email}, "Staymate. Confirm email", fmt.Sprintf("Hello it's StayMate. Verify your email.\nCode: %v", code), "", "", []byte{})
	if err != nil {
		myLogger.Log.Error().Err(err).Msg("Error sending verification code email")
		http.Error(w, "Failed to send verification code", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Verification code sent successfully"))
}

// Verify handles the email verification process
func (h *verifyHandler) Verify(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming JSON body
	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		myLogger.Log.Error().Err(err).Msg("Error decoding JSON request body")
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

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

	// Check if verification code expired
	if time.Now().After(userEmailConfirm.ExpiresAt) {
		myLogger.Log.Warn().Str("email", request.Email).Msg("Verification code has expired")
		http.Error(w, "Verification code has expired", http.StatusBadRequest)
		return
	}

	// If email and code match, update the status of the user to "active"
	var user models.User
	err = h.db.Where("email = ?", request.Email).First(&user).Error
	if err != nil {
		myLogger.Log.Error().Err(err).Str("email", request.Email).Msg("Failed to find user in Users table")
		http.Error(w, "User with this email does not exist", http.StatusBadRequest)
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

	w.Write([]byte("Email successfully verified, user status is now active"))
}

// Generate a 4-digit code between 1000 and 9999
func GenerateCode() string {
	code := rand.IntN(9000) + 1000
	return fmt.Sprintf("%d", code)
}
