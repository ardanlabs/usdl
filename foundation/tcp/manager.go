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
func (cm *ClientManager) Shutdown(ctx context.Context) error {
	cm.log(EvtStop, TypInfo, "", "client manager started shutdown")
	defer cm.log(EvtStop, TypInfo, "", "client manager completed shutdown")

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()

		for _, c := range cm.clients.copy() {
			go c.close()
		}
	}()

	<-ctx.Done()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		cm.log(EvtStop, TypInfo, "", "client manager deadline exceeded")
		return ctx.Err()
	}

	return nil
}

// Dial establishes a new TCP connection to the specified address.
func (cm *ClientManager) Dial(network string, address string) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	// Add this new connection to the manager map and
	// start the client goroutine.
	clt, err := cm.startNewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("startNewClient: %w", err)
	}

	return clt, nil
}

// =============================================================================

// startNewClient takes a new connection and adds it to the manager.
func (cm *ClientManager) startNewClient(conn net.Conn) (*Client, error) {
	tcpAddr := conn.LocalAddr().(*net.TCPAddr)

	if _, err := cm.clients.find(tcpAddr); err == nil {
		cm.log(EvtJoin, TypError, tcpAddr.IP.String(), "already connected")
		conn.Close()
		return nil, errors.New("client already connected")
	}

	clt := newClient(cm.log, cm.clients, cm.handlers, conn)

	cm.clients.add(clt)

	clt.start()

	return clt, nil
}
