package handler

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type ChatWebsocketHandler struct {
	upgrader    websocket.Upgrader
	connections map[string]*websocket.Conn // userID -> WebSocket
	adminConn   *websocket.Conn            // Admin WebSocket
	sync.Mutex

	logger *zerolog.Logger
}

func NewChatWebsocketHandler(logger *zerolog.Logger) *ChatWebsocketHandler {
	return &ChatWebsocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		connections: make(map[string]*websocket.Conn),
		logger:      logger,
	}
}

// WebSocket handler for users
func (h *ChatWebsocketHandler) UserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	conn, err := h.upgrader.Upgrade(w, r, nil)

	if err != nil {
		h.logger.Error().Err(err).Msg("User websocket connection error")
		http.Error(w, "WebSocket error", http.StatusInternalServerError)
		return
	}

	h.Lock()
	h.connections[userID] = conn
	h.Unlock()

	h.logger.Info().Str("userID", userID).Msg("User connected to chat")

	defer func() {
		h.Lock()
		delete(h.connections, userID)
		h.Unlock()
		conn.Close()
		h.logger.Warn().Str("userID", userID).Msg("User disconnected from chat")
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			h.logger.Warn().Str("userID", userID).Err(err).Msg("Error receiving message from user")
			break
		}

		h.logger.Info().Str("userID", userID).Str("message", string(msg)).Msg("Message from user")

		// Send message to admin
		h.Lock()
		if h.adminConn != nil {
			err = h.adminConn.WriteMessage(websocket.TextMessage, []byte(userID+": "+string(msg)))
			if err != nil {
				h.logger.Error().Err(err).Msg("Error sending message to admin")
			}
		} else {
			h.logger.Warn().Str("userID", userID).Msg("Message sent, but admin is not connected")
		}
		h.Unlock()
	}
}

// WebSocket handler for admin
func (h *ChatWebsocketHandler) AdminHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error().Err(err).Msg("Admin connection error")
		http.Error(w, "WebSocket error", http.StatusInternalServerError)
		return
	}

	h.Lock()
	h.adminConn = conn
	h.Unlock()

	h.logger.Info().Msg("Admin connected to chat")

	defer func() {
		h.Lock()
		h.adminConn = nil
		h.Unlock()
		conn.Close()
		h.logger.Warn().Msg("Admin disconnected from chat")
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			h.logger.Warn().Err(err).Msg("Error receiving message from admin")
			break
		}

		// Expecting format "userID: message"
		// Expecting format "userID: message"
		msgParts := strings.SplitN(string(msg), ": ", 2) // Split only at the first ": "
		if len(msgParts) < 2 {
			h.logger.Warn().Str("rawMessage", string(msg)).Msg("Invalid message format from admin")
			return
		}
		userID, text := msgParts[0], msgParts[1]

		h.logger.Info().Str("userID", userID).Str("message", text).Msg("Message from admin")

		h.Lock()
		if userConn, exists := h.connections[userID]; exists {
			err = userConn.WriteMessage(websocket.TextMessage, []byte("Admin: "+text))
			if err != nil {
				h.logger.Error().Err(err).Str("userID", userID).Msg("Error sending message to user")
			}
		} else {
			h.logger.Warn().Str("userID", userID).Msg("Admin sent a message, but user is not connected")
		}
		h.Unlock()
	}
}
