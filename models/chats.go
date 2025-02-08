package models

import (
	"time"

	"github.com/google/uuid"
)

// Chat represents a chat session
type Chat struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ChatUUID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"chat_uuid"` // Уникальный UUID чата
	UserID    uint      `gorm:"not null;index" json:"user_id"`                   // Внешний ключ -> users.id
	Status    string    `gorm:"type:varchar(10);default:'active'" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Messages []Message `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"messages"`
}

// Message represents a chat message
type Message struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ChatID    uint      `gorm:"not null;index" json:"chat_id"` // Foreign Key -> chats.id
	Sender    string    `gorm:"type:varchar(10);not null;check:sender IN ('user', 'admin')" json:"sender"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Chat Chat `gorm:"foreignKey:ChatID" json:"-"`
}
