package chat

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// User represents a user in the chat system.
type User struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	LastPing time.Time       `json:"lastPing"`
	LastPong time.Time       `json:"lastPong"`
	Conn     *websocket.Conn `json:"-"`
}

// Connection represents a connection to a user.
type Connection struct {
	Conn     *websocket.Conn
	LastPing time.Time
	LastPong time.Time
}

type inMessage struct {
	ToID string `json:"toID"`
	Msg  string `json:"msg"`
}

type outMessage struct {
	From User   `json:"from"`
	Msg  string `json:"msg"`
}

type busMessage struct {
	CapID    uuid.UUID `json:"capID"`
	FromID   string    `json:"fromID"`
	FromName string    `json:"fromName"`
	ToID     string    `json:"toID"`
	Msg      string    `json:"msg"`
}
