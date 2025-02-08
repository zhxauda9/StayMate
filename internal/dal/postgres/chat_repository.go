package postgres

import (
	"errors"

	"github.com/google/uuid"
	"github.com/zhxauda9/StayMate/models"
	"gorm.io/gorm"
)

// ChatRepository - интерфейс репозитория
type ChatRepository interface {
	CreateChat(userID uint) (*models.Chat, error)
	GetChatByUUID(chatUUID uuid.UUID) (*models.Chat, error)
	CloseChat(chatUUID uuid.UUID) error
	GetActiveChats() ([]models.Chat, error)
	SaveMessage(chatUUID uuid.UUID, sender, message string) error
	GetMessages(chatUUID uuid.UUID) ([]models.Message, error)
}

// ChatRepo - реализация репозитория
type ChatRepo struct {
	db *gorm.DB
}

// NewChatRepository - конструктор репозитория
func NewChatRepository(db *gorm.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

// CreateChat - создаёт новый чат
func (r *ChatRepo) CreateChat(userID uint) (*models.Chat, error) {
	newChat := models.Chat{
		ChatUUID: uuid.New(),
		UserID:   userID,
		Status:   "active",
	}

	if err := r.db.Create(&newChat).Error; err != nil {
		return nil, err
	}
	return &newChat, nil
}

// GetChatByID - get Chat By ID
func (r *ChatRepo) GetChatByUUID(chatUUID uuid.UUID) (*models.Chat, error) {
	var chat models.Chat
	if err := r.db.Where("chat_uuid = ?", chatUUID).Preload("Messages").First(&chat).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

// CloseChat - Makes chat incative
func (r *ChatRepo) CloseChat(chatUUID uuid.UUID) error {
	return r.db.Model(&models.Chat{}).Where("chat_uuid = ?", chatUUID).Update("status", "inactive").Error
}

// GetActiveChats - Gets all active chats
func (r *ChatRepo) GetActiveChats() ([]models.Chat, error) {
	var chats []models.Chat
	if err := r.db.Where("status = ?", "active").Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

// SaveMessage - Saves message
func (r *ChatRepo) SaveMessage(chatUUID uuid.UUID, sender, message string) error {
	var chat models.Chat
	err := r.db.Where("chat_uuid = ?", chatUUID).First(&chat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			chat = models.Chat{ChatUUID: chatUUID}
			if err := r.db.Create(&chat).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	msg := models.Message{
		ChatID:  chat.ID,
		Sender:  sender,
		Message: message,
	}

	return r.db.Create(&msg).Error
}

// GetMessages - Get messages from chat by chatID
func (r *ChatRepo) GetMessages(chatUUID uuid.UUID) ([]models.Message, error) {
	var chat models.Chat
	if err := r.db.Where("chat_uuid = ?", chatUUID).First(&chat).Error; err != nil {
		return nil, err
	}

	var messages []models.Message
	if err := r.db.Where("chat_id = ?", chat.ID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
