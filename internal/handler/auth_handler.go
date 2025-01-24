// internal/handler/auth_handler.go
package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/zhxauda9/StayMate/internal/service"
	"github.com/zhxauda9/StayMate/models"
)

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *authHandler {
	return &authHandler{authService: authService}
}

func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	verificationCode := generateVerificationCode()
	user.VerificationCode = verificationCode

	if err := h.authService.Register(user); err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		fmt.Errorf("failed to create user: %w", err)
		fmt.Println("failed to create user ", err)
		return
	}

	sendEmail(user.Email, verificationCode)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, fmt.Sprintf("Invalid input: %v", err), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid credentials: %v", err), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600 * 24 * 30,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful", "token": token})
}

func (h *authHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value

	isValid, err := h.authService.ValidateToken(tokenString)
	if err != nil || !isValid {
		http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "I'm logging!", "token": tokenString})
}

func (h *authHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}

	tokenString := cookie.Value
	isValid, err := h.authService.ValidateToken(tokenString)
	if err != nil || !isValid {
		http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
		return
	}
	user, err := h.authService.GetUserFromToken(tokenString)
	if err != nil {
		http.Error(w, "Failed to get user data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   -1,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}

func (h *authHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var request struct {
		VerificationCode string `json:"verificationCode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := h.authService.FindUserByVerificationCode(request.VerificationCode)
	if err != nil || user.Verified {
		http.Error(w, "Invalid or expired code", http.StatusUnauthorized)
		return
	}

	user.Verified = true
	user.VerificationCode = ""
	if err := h.authService.UpdateUser(user.ID, user); err != nil {
		http.Error(w, "Failed to verify user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User verified successfully"})
}

func sendEmail(emailToCheck string, verificationCode string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	mailServ, err := service.NewMailService(smtpHost, smtpPort, email, password)
	subject := "Email Verification"
	message := fmt.Sprintf("Your verification code is: %s", verificationCode)

	err = mailServ.Send([]string{emailToCheck}, subject, message, "", "", nil)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
