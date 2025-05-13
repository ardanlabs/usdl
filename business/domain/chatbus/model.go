package chatbus

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// User represents a user in the chat system.
type User struct {
	ID       common.Address  `json:"id"`
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

type incomingMessage struct {
	ToID      common.Address `json:"toID"`
	Encrypted bool           `json:"encrypted"`
	Msg       [][]byte       `json:"msg"`
	FromNonce uint64         `json:"fromNonce"`
	V         *big.Int       `json:"v"`
	R         *big.Int       `json:"r"`
	S         *big.Int       `json:"s"`
}

type outgoingUser struct {
	ID    common.Address `json:"id"`
	Name  string         `json:"name"`
	Nonce uint64         `json:"nonce"`
}

type outgoingMessage struct {
	From      outgoingUser `json:"from"`
	Encrypted bool         `json:"encrypted"`
	Msg       [][]byte     `json:"msg"`
}

type busMessage struct {
	CapID    uuid.UUID      `json:"capID"`
	FromID   common.Address `json:"fromID"`
	FromName string         `json:"fromName"`
	incomingMessage
}
