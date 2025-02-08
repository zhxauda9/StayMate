package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/zhxauda9/StayMate/internal/dal/postgres"
	"github.com/zhxauda9/StayMate/models"
)

type ChatHandler struct {
	chatRepo postgres.ChatRepository
}

func NewChatHandler(chatRepo postgres.ChatRepository) *ChatHandler {
	return &ChatHandler{chatRepo: chatRepo}
}

// takes UUID in path and return page with chat by given UUID
func (h *ChatHandler) AdminChatPage(w http.ResponseWriter, r *http.Request) {
	chatUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid chat UUID", http.StatusBadRequest)
		return
	}

	// Получаем чат и сообщения из БД
	chat, err := h.chatRepo.GetChatByUUID(chatUUID)
	if err != nil {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}

	dataB, _ := json.MarshalIndent(chat, " ", "    ")
	fmt.Println(string(dataB))
	tmpl, err := template.ParseFiles("web/templates/admin-chat.html")
	if err != nil {
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

	tmpl.Execute(w, data)
}

// StartChat - creates new chat
func (h *ChatHandler) StartChat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID uint `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	chat, err := h.chatRepo.CreateChat(req.UserID)
	if err != nil {
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Chat started",
		"chat_id": chat.ID,
	})
}

// GetChatHistory - получает историю сообщений
func (h *ChatHandler) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	chatUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid chat UUID", http.StatusBadRequest)
		return
	}

	chat, err := h.chatRepo.GetChatByUUID(chatUUID)
	if err != nil {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

// CloseChat - закрывает чат
func (h *ChatHandler) CloseChat(w http.ResponseWriter, r *http.Request) {
	chatUUID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid chat UUID", http.StatusBadRequest)
		return
	}

	if err := h.chatRepo.CloseChat(chatUUID); err != nil {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Chat closed"})
}

// GetActiveChats - получает все активные чаты
func (h *ChatHandler) GetActiveChats(w http.ResponseWriter, r *http.Request) {
	chats, err := h.chatRepo.GetActiveChats()
	if err != nil {
		http.Error(w, "Failed to fetch chats", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}

// func (h *ChatHandler) GetMessagesByChatUUID(w http.ResponseWriter, r *http.Request) {
// 	// Получаем UUID из URL
// 	chatUUID, err := uuid.Parse(r.PathValue("id"))
// 	if err != nil {
// 		http.Error(w, `{"error": "Invalid UUID format"}`, http.StatusBadRequest)
// 		return
// 	}

// 	// Получаем сообщения чата
// 	messages, err := h.chatRepo.GetMessages(chatUUID)
// 	if err != nil {
// 		http.Error(w, `{"error": "Failed to retrieve messages"}`, http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(messages)
// }
