package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"
)

// ClientConfig provides a data structure of required configuration parameters.
type ClientConfig struct {
	Handlers Handlers // Support for binding and handling requests.
	Logger   Logger   // Support for logging events that occur in the TCP listener.
}

func (cfg ClientConfig) validate() error {
	if cfg.Handlers == nil {
		return ErrInvalidHandlers
	}

	if cfg.Logger == nil {
		return ErrInvalidLoggerHandler
	}

	return nil
}

// ClientManager manages a collection of TCP client connections.
type ClientManager struct {
	log      internalLogger
	handlers Handlers
	clients  *clients
}

// NewClientManager creates a new ClientManager.
func NewClientManager(cfg ClientConfig) (*ClientManager, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	l := func(evt int, typ int, ipAddress string, format string, a ...any) {
		cfg.Logger(eventTypes[evt], eventSubTypes[typ], ipAddress, fmt.Sprintf(format, a...))
	}

	cm := ClientManager{
		log:      l,
		handlers: cfg.Handlers,
		clients:  newClients(l),
	}

	return &cm, nil
}

// Shutdown shuts down the manager and closes all connections.
func (cln *ClientManager) Shutdown(ctx context.Context) error {
	cln.log(EvtStop, TypInfo, "", "started shutdown")
	defer cln.log(EvtStop, TypInfo, "", "completed shutdown")

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()

		for _, c := range cln.clients.copy() {
			go c.close()
		}
	}()

	<-ctx.Done()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		cln.log(EvtStop, TypInfo, "", "deadline exceeded")
		return ctx.Err()
	}

	return nil
}

// Dial establishes a new TCP connection to the specified address.
func (cln *ClientManager) Dial(network string, address string) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	// Add this new connection to the manager map and
	// start the client goroutine.
	cln.startNewClient(conn)

	return nil
}

// =============================================================================

// startNewClient takes a new connection and adds it to the manager.
func (cln *ClientManager) startNewClient(conn net.Conn) {
	tcpAddr := conn.RemoteAddr().(*net.TCPAddr)

	if _, err := cln.clients.find(tcpAddr); err == nil {
		cln.log(EvtJoin, TypError, tcpAddr.IP.String(), "already connected")
		conn.Close()
		return
	}

	c := newClient(cln.log, cln.clients, cln.handlers, conn)

	cln.clients.add(c)

	c.start()
}
