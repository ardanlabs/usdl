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
	name     string
	log      internalLogger
	handlers Handlers
	clients  *clients
}

// NewClientManager creates a new ClientManager.
func NewClientManager(name string, cfg ClientConfig) (*ClientManager, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	l := func(ctx context.Context, name string, evt int, typ int, ipAddress string, format string, a ...any) {
		cfg.Logger(ctx, name, eventTypes[evt], eventSubTypes[typ], ipAddress, fmt.Sprintf(format, a...))
	}

	cm := ClientManager{
		name:     name,
		log:      l,
		handlers: cfg.Handlers,
		clients:  newClients(l),
	}

	return &cm, nil
}

// Shutdown shuts down the manager and closes all connections.
func (cm *ClientManager) Shutdown(ctx context.Context) error {
	cm.log(ctx, cm.name, EvtStop, TypInfo, "", "client manager started shutdown")
	defer cm.log(ctx, cm.name, EvtStop, TypInfo, "", "client manager completed shutdown")

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()

		for _, c := range cm.clients.copy() {
			go c.close()
		}
	}()

	<-ctx.Done()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		cm.log(ctx, cm.name, EvtStop, TypInfo, "", "client manager deadline exceeded")
		return ctx.Err()
	}

	return nil
}

// Dial establishes a new TCP connection to the specified address.
func (cm *ClientManager) Dial(ctx context.Context, key string, network string, address string) (*Client, error) {
	if _, err := cm.clients.find(key); err == nil {
		return nil, errors.New("client already connected")
	}

	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	clt, err := newClient(key, cm.name, cm.log, cm.clients, cm.handlers, conn)
	if err != nil {
		return nil, fmt.Errorf("newClient: %w", err)
	}

	cm.clients.add(key, clt)

	clt.start()

	return clt, nil
}

// Retrieve retrieves a client by user ID.
func (cm *ClientManager) Retrieve(ctx context.Context, key string) (*Client, error) {
	clt, err := cm.clients.find(key)
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	return clt, nil
}
