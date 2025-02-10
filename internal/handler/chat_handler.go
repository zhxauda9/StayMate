package handler

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/zhxauda9/StayMate/internal/dal/postgres"
	"github.com/zhxauda9/StayMate/models"
)

type UserByEmailInterface interface {
	GetUserByEmail(email string) (models.User, error)
}

type ChatHandler struct {
	chatRepo postgres.ChatRepository
	userRepo UserByEmailInterface
	logger   *zerolog.Logger
}

func NewChatHandler(chatRepo postgres.ChatRepository, userRepo UserByEmailInterface, logger *zerolog.Logger) *ChatHandler {
	return &ChatHandler{chatRepo: chatRepo, userRepo: userRepo, logger: logger}
}

// Renders admin page using templates
func (h *ChatHandler) AdminChatPage(w http.ResponseWriter, r *http.Request) {
	chatUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid chat UUID")
		http.Error(w, "Invalid chat UUID", http.StatusBadRequest)
		return
	}

	h.logger.Info().Str("chat_uuid", chatUUID.String()).Msg("Fetching chat page")

	chat, err := h.chatRepo.GetChatByUUID(chatUUID)
	if err != nil {
		h.logger.Error().Err(err).Str("chat_uuid", chatUUID.String()).Msg("Chat not found")
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/admin-chat.html")
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to load template")
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	data := struct {
		ChatUUID string
		Messages []models.Message
	}{
		ChatUUID: chat.ChatUUID.String(),
		Messages: chat.Messages,
	}

	h.logger.Info().Str("chat_uuid", chatUUID.String()).Msg("Rendering chat page")
	tmpl.Execute(w, data)
}

// Starts new chat for user and admin
func (h *ChatHandler) StartChat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Info().Str("email", req.Email).Msg("Starting chat")

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		h.logger.Error().Err(err).Str("email", req.Email).Msg("User not found")
		http.Error(w, "Email does not exist", http.StatusBadRequest)
		return
	}

	chat, err := h.chatRepo.CreateChat(uint(user.ID))
	if err != nil {
		h.logger.Error().Err(err).Str("email", req.Email).Msg("Failed to create chat")
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	h.logger.Info().Str("chat_uuid", chat.ChatUUID.String()).Msg("Chat created successfully")

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_chat_uuid",
		Value:    chat.ChatUUID.String(),
		Path:     "/",
		HttpOnly: false,                // Allow JavaScript access
		Secure:   false,                // Secure must be false for HTTP (localhost)
		SameSite: http.SameSiteLaxMode, // Prevents cross-site issues while allowing normal usage
		MaxAge:   3600 * 24 * 7,        // 1 week
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Chat started",
		"chat_uuid": chat.ChatUUID.String(),
	})
}

// Return information and history of the chat
func (h *ChatHandler) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	chatUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid chat UUID")
		http.Error(w, "Invalid chat UUID", http.StatusBadRequest)
		return
	}

	h.logger.Info().Str("chat_uuid", chatUUID.String()).Msg("Fetching chat history")

	chat, err := h.chatRepo.GetChatByUUID(chatUUID)
	if err != nil {
		h.logger.Error().Err(err).Str("chat_uuid", chatUUID.String()).Msg("Chat not found")
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

// Closes the chat
func (h *ChatHandler) CloseChat(w http.ResponseWriter, r *http.Request) {
	chatUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid chat UUID")
		http.Error(w, "Invalid chat UUID", http.StatusBadRequest)
		return
	}

	h.logger.Info().Str("chat_uuid", chatUUID.String()).Msg("Closing chat")

	if err := h.chatRepo.CloseChat(chatUUID); err != nil {
		h.logger.Error().Err(err).Str("chat_uuid", chatUUID.String()).Msg("Chat not found or already closed")
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Chat closed"})
}

// Returns chats with status "active"
func (h *ChatHandler) GetActiveChats(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("Fetching active chats")

	chats, err := h.chatRepo.GetActiveChats()
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to fetch active chats")
		http.Error(w, "Failed to fetch chats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}
