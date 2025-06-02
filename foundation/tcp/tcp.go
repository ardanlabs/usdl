package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type internalLogger func(evt int, typ int, ipAddress string, format string, a ...any)

// TCP contains a set of networked client connections.
type TCP struct {
	name                   string
	log                    internalLogger
	netType                string
	addr                   string
	connHandler            ConnHandler
	reqHandler             ReqHandler
	respHandler            RespHandler
	ipAddress              string
	port                   int
	tcpAddr                *net.TCPAddr
	listener               *listener
	clients                *clients
	wgStartG               sync.WaitGroup
	shuttingDown           atomic.Bool
	lastAcceptedConnection time.Time
}

// New creates a new manager to service clients.
func New(name string, cfg Config) (*TCP, error) {

	// Validate the configuration.
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	// Resolve the addr that is provided.
	tcpAddr, err := net.ResolveTCPAddr(cfg.NetType, cfg.Addr)
	if err != nil {
		return nil, err
	}

	l := func(evt int, typ int, ipAddress string, format string, a ...any) {
		cfg.Logger(eventTypes[evt], eventSubTypes[typ], ipAddress, fmt.Sprintf(format, a...))
	}

	// Create a TCP for this ipaddress and port.
	t := TCP{
		name:        name,
		log:         l,
		netType:     cfg.NetType,
		addr:        cfg.Addr,
		connHandler: cfg.ConnHandler,
		reqHandler:  cfg.ReqHandler,
		respHandler: cfg.RespHandler,
		ipAddress:   tcpAddr.IP.String(),
		port:        tcpAddr.Port,
		tcpAddr:     tcpAddr,
		listener:    newListener(),
		clients:     newClients(l),
	}

	return &t, nil
}

// Name returns the name of the TCP manager.
func (t *TCP) Name() string {
	return t.name
}

// Shutdown shuts down the manager and closes all connections.
func (t *TCP) Shutdown(ctx context.Context) error {
	t.log(EvtStop, TypInfo, "", "started shutdown")
	defer t.log(EvtStop, TypInfo, "", "completed shutdown")

	t.shuttingDown.Store(true)

	t.listener.reset()

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()

		for _, c := range t.clients.copy() {
			go c.close()
		}
	}()

	<-ctx.Done()

	if ctx.Err() != nil {
		t.log(EvtStop, TypInfo, "", "cancelled shutdown")
		return ctx.Err()
	}

	t.wgStartG.Wait()

	return nil
}

// Listen creates the accept routine and begins to accept connections.
func (t *TCP) Listen() error {
	if t.listener.tcpListener() != nil {
		return errors.New("this TCP has already been started")
	}

	t.wgStartG.Add(1)

	go func() {
		defer func() {
			t.log(EvtAccept, TypInfo, net.JoinHostPort(t.ipAddress, strconv.Itoa(t.port)), "completed listener shutdown")
			t.wgStartG.Done()
		}()

	startlistener:
		for {
			if t.shuttingDown.Load() {
				t.log(EvtAccept, TypInfo, net.JoinHostPort(t.ipAddress, strconv.Itoa(t.port)), "started listener shutdown")
				t.listener.reset()
				break
			}

			listener, err := t.listener.start(t.netType, t.tcpAddr)
			if err != nil {
				// TODO: Use Context to control the retry / cancel.
				t.log(EvtAccept, TypError, "", err.Error())
				time.Sleep(200 * time.Millisecond)
				continue
			}

			t.log(EvtAccept, TypInfo, net.JoinHostPort(t.ipAddress, strconv.Itoa(t.port)), "waiting")

			for {
				conn, err := listener.Accept()
				if err != nil {
					if t.shuttingDown.Load() {
						t.log(EvtAccept, TypInfo, net.JoinHostPort(t.ipAddress, strconv.Itoa(t.port)), "started listener shutdown")
						t.listener.reset()
						break startlistener
					}

					t.log(EvtAccept, TypError, conn.RemoteAddr().String(), err.Error())

					type temporary interface {
						Temporary() bool
					}

					if e, ok := err.(temporary); ok && !e.Temporary() {
						t.listener.reset()
						continue startlistener
					}

					continue
				}

				// Add this new connection to the manager map and
				// start the client goroutine.
				t.startNewClient(conn)
			}
		}
	}()

	t.wgStartG.Wait()

	return nil
}

// CloseClient will close the client socket connection.
func (t *TCP) CloseClient(tcpAddr *net.TCPAddr) error {
	c, err := t.clients.find(tcpAddr)
	if err != nil {
		return fmt.Errorf("IP[ %s ] : disconnected", tcpAddr.String())
	}

	// Drop the connection using a goroutine since we are on the
	// socket goroutine most likely.
	go c.close()

	return nil
}

// Send will deliver the response back to the client.
func (t *TCP) Send(ctx context.Context, r *Response) error {
	c, err := t.clients.find(r.TCPAddr)
	if err != nil {
		return fmt.Errorf("IP[ %s ] : disconnected", r.TCPAddr.String())
	}

	return t.respHandler.Write(r, c.writer)
}

// SendAll will deliver the response back to all connected clients.
func (t *TCP) SendAll(ctx context.Context, r *Response) error {
	clients := t.clients.copy()

	var errors Errors
	for _, c := range clients {
		if err := t.respHandler.Write(r, c.writer); err != nil {
			errors = append(errors, err)
		}
	}

	if errors != nil {
		return errors
	}

	return nil
}

// Addr returns the listener's network address. This may be different than the values
// provided in the configuration, for example if configuration port value is 0.
func (t *TCP) Addr() net.Addr {
	return t.tcpAddr
}

// Clients returns the number of active clients connected.
func (t *TCP) Clients() int {
	return t.clients.count()
}

// Groom drops connections that are not active for the specified duration.
func (t *TCP) Groom(d time.Duration) {
	client := t.clients.copy()

	now := time.Now().UTC()
	for _, c := range client {
		sub := now.Sub(c.lastAct)
		if sub >= d {
			// TODO
			// This is a blocking call that waits for the socket goroutine
			// to report its done. This parallel call should work well since
			// there is no error handling needed.
			t.log(EvtGroom, TypInfo, c.ipAddress, "Last[ %v ] Dur[ %v ]", c.lastAct.Format(time.RFC3339), sub)
			go c.close()
		}
	}
}

// =============================================================================

// startNewClient takes a new connection and adds it to the manager.
func (t *TCP) startNewClient(conn net.Conn) {
	tcpAddr := conn.RemoteAddr().(*net.TCPAddr)

	if _, err := t.clients.find(tcpAddr); err == nil {
		t.log(EvtJoin, TypError, tcpAddr.IP.String(), "already connected")
		conn.Close()
		return
	}

	c := newClient(t, conn)

	t.clients.add(c)

	c.start()
}
