// Package tcpapp provides the application layer for the TCP domain.
package tcpapp

import (
	"context"
	"fmt"

	"github.com/ardanlabs/usdl/foundation/logger"
	"github.com/ardanlabs/usdl/foundation/tcp"
)

// ServerHandlers implements the Handlers interface for the TCP server.
type ServerHandlers struct {
	log *logger.Logger
}

// NewServerHandlers creates a new instance of ServerHandlers.
func NewServerHandlers(log *logger.Logger) *ServerHandlers {
	return &ServerHandlers{
		log: log,
	}
}

// Bind binds the client to the server handlers.
func (sh ServerHandlers) Bind(ctx context.Context, clt *tcp.Client) {
	sh.log.Info(ctx, "bind")
}

// Read reads data from the client connection.
func (ServerHandlers) Read(ctx context.Context, clt *tcp.Client) ([]byte, int, error) {
	return nil, 0, fmt.Errorf("not implemented")
}

// Process processes the request from the client.
func (ServerHandlers) Process(ctx context.Context, r *tcp.Request, clt *tcp.Client) {

}
