package chatapp

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

type tcpConnRequest struct {
	TUIUserID    string `json:"tui_user_id"`
	ClientUserID string `json:"client_user_id"`
	TCPHost      string `json:"tcp_host"`
}

// Decode implements the decoder interface.
func (app *tcpConnRequest) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

type tcpConnDropResponse struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message"`
}

func (app tcpConnDropResponse) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

type stateResponse struct {
	TCPConnections []common.Address `json:"tcp_connections"`
}

func (app stateResponse) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}
