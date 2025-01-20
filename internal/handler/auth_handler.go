// internal/handler/auth_handler.go
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zhxauda9/StayMate/internal/service"
	"github.com/zhxauda9/StayMate/models"
)

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *authHandler {
	return &authHandler{authService: authService}
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.authService.Register(user); err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

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
