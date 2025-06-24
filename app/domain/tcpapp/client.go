// Package tcpapp provides the application layer for the TCP domain.
package tcpapp

import (
	"fmt"

	"github.com/ardanlabs/usdl/foundation/logger"
	"github.com/ardanlabs/usdl/foundation/tcp"
)

// ClientHandlers implements the Handlers interface for the TCP client manager.
type ClientHandlers struct {
	log *logger.Logger
}

// NewClientHandlers creates a new instance of ClientHandlers.
func NewClientHandlers(log *logger.Logger) *ClientHandlers {
	return &ClientHandlers{
		log: log,
	}
}

// Bind binds the client to the server handlers.
func (ch ClientHandlers) Bind(clt *tcp.Client) {
	ch.log.Info(clt.Context(), "bind")
}

// Read reads data from the client connection.
func (ClientHandlers) Read(clt *tcp.Client) ([]byte, int, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

// Process processes the request from the client.
func (ClientHandlers) Process(r *tcp.Request, clt *tcp.Client) {
}

// Drop is called when a connection is dropped.
func (ch ClientHandlers) Drop(clt *tcp.Client) {
	ch.log.Info(clt.Context(), "drop")
}
