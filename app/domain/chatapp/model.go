package chatapp

import "encoding/json"

type tcpConnRequest struct {
	UserID  string `json:"user_id"`
	TCPHost string `json:"tcp_host"`
}

// Decode implements the decoder interface.
func (app *tcpConnRequest) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}
