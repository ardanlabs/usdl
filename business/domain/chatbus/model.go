package chatbus

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// UIUser represents a web socket user in the chat system.
type UIUser struct {
	ID       common.Address  `json:"id"`
	Name     string          `json:"name"`
	LastPing time.Time       `json:"lastPing"`
	LastPong time.Time       `json:"lastPong"`
	UIConn   *websocket.Conn `json:"-"`
}

// UIConnection represents a connection to a user.
type UIConnection struct {
	Conn     *websocket.Conn
	LastPing time.Time
	LastPong time.Time
}

type uiIncomingMessage struct {
	ToID      common.Address `json:"toID"`
	Encrypted bool           `json:"encrypted"`
	Msg       [][]byte       `json:"msg"`
	FromNonce uint64         `json:"fromNonce"`
	V         *big.Int       `json:"v"`
	R         *big.Int       `json:"r"`
	S         *big.Int       `json:"s"`
}

type uiOutgoingUser struct {
	ID    common.Address `json:"id"`
	Name  string         `json:"name"`
	Nonce uint64         `json:"nonce"`
}

type uiOutgoingMessage struct {
	From      uiOutgoingUser `json:"from"`
	Encrypted bool           `json:"encrypted"`
	Msg       [][]byte       `json:"msg"`
}

type natsInOutMessage struct {
	CapID    uuid.UUID      `json:"capID"`
	FromID   common.Address `json:"fromID"`
	FromName string         `json:"fromName"`
	uiIncomingMessage
}
