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

	"github.com/google/uuid"
)

// ServerConfig provides a data structure of required configuration parameters.
type ServerConfig struct {
	NetType  string   // "tcp", tcp4" or "tcp6"
	Addr     string   // "host:port" or "[ipv6-host%zone]:port"
	Handlers Handlers // Support for binding and handling requests.
	Logger   Logger   // Support for logging events that occur in the TCP listener.
}

func (cfg ServerConfig) validate() error {
	if cfg.NetType != "tcp" && cfg.NetType != "tcp4" && cfg.NetType != "tcp6" {
		return ErrInvalidNetType
	}

	if cfg.Handlers == nil {
		return ErrInvalidHandlers
	}

	if cfg.Logger == nil {
		return ErrInvalidLoggerHandler
	}

	return nil
}

// Server contains a set of networked client connections.
type Server struct {
	ctx                    context.Context
	name                   string
	log                    internalLogger
	netType                string
	addr                   string
	handlers               Handlers
	ipAddress              string
	port                   int
	tcpAddr                *net.TCPAddr
	listener               *listener
	clients                *clients
	wgStartG               sync.WaitGroup
	shuttingDown           atomic.Bool
	lastAcceptedConnection time.Time
}

// NewServer creates an API for a TCP server that can accept connections.
func NewServer(name string, cfg ServerConfig) (*Server, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	tcpAddr, err := net.ResolveTCPAddr(cfg.NetType, cfg.Addr)
	if err != nil {
		return nil, err
	}

	l := func(ctx context.Context, name string, evt int, typ int, ipAddress string, format string, a ...any) {
		cfg.Logger(ctx, name, eventTypes[evt], eventSubTypes[typ], ipAddress, fmt.Sprintf(format, a...))
	}

	t := Server{
		ctx:       setTraceID(context.Background(), uuid.New()),
		name:      name,
		log:       l,
		netType:   cfg.NetType,
		addr:      cfg.Addr,
		handlers:  cfg.Handlers,
		ipAddress: tcpAddr.IP.String(),
		port:      tcpAddr.Port,
		tcpAddr:   tcpAddr,
		listener:  newListener(),
		clients:   newClients(l),
	}

	return &t, nil
}

// Shutdown shuts down the manager and closes all connections.
func (srv *Server) Shutdown(ctx context.Context) error {
	srv.log(ctx, srv.name, EvtStop, TypInfo, "", "server started shutdown")
	defer srv.log(ctx, srv.name, EvtStop, TypInfo, "", "server completed shutdown")

	srv.shuttingDown.Store(true)

	srv.listener.reset()

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()

		for _, c := range srv.clients.copy() {
			go c.close()
		}
	}()

	<-ctx.Done()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		srv.log(ctx, srv.name, EvtStop, TypInfo, "", "server deadline exceeded")
		return ctx.Err()
	}

	srv.wgStartG.Wait()

	return nil
}

// Name returns the name of the TCP manager.
func (srv *Server) Name() string {
	return srv.name
}

// Listen creates the accept routine and begins to accept connections.
func (srv *Server) Listen() error {
	if srv.listener.tcpListener() != nil {
		return errors.New("this TCP has already been started")
	}

	srv.wgStartG.Add(1)

	go func() {
		defer func() {
			srv.log(srv.ctx, srv.name, EvtAccept, TypInfo, net.JoinHostPort(srv.ipAddress, strconv.Itoa(srv.port)), "completed listener shutdown")
			srv.wgStartG.Done()
		}()

	startlistener:
		for {
			if srv.shuttingDown.Load() {
				srv.log(srv.ctx, srv.name, EvtAccept, TypInfo, net.JoinHostPort(srv.ipAddress, strconv.Itoa(srv.port)), "started listener shutdown")
				srv.listener.reset()
				break
			}

			listener, err := srv.listener.start(srv.netType, srv.tcpAddr)
			if err != nil {
				// TODO: Use Context to control the retry / cancel.
				srv.log(srv.ctx, srv.name, EvtAccept, TypError, "", err.Error())
				time.Sleep(200 * time.Millisecond)
				continue
			}

			for {
				srv.log(srv.ctx, srv.name, EvtAccept, TypInfo, net.JoinHostPort(srv.ipAddress, strconv.Itoa(srv.port)), "waiting")

				conn, err := listener.Accept()
				if err != nil {
					if srv.shuttingDown.Load() {
						srv.log(srv.ctx, srv.name, EvtAccept, TypInfo, net.JoinHostPort(srv.ipAddress, strconv.Itoa(srv.port)), "started listener shutdown")
						srv.listener.reset()
						break startlistener
					}

					srv.log(srv.ctx, srv.name, EvtAccept, TypError, conn.RemoteAddr().String(), err.Error())

					type temporary interface {
						Temporary() bool
					}

					if e, ok := err.(temporary); ok && !e.Temporary() {
						srv.listener.reset()
						continue startlistener
					}

					continue
				}

				// Add this new connection to the manager map and
				// start the client goroutine.
				if err := srv.startNewClient(conn); err != nil {
					srv.log(srv.ctx, srv.name, EvtAccept, TypError, conn.RemoteAddr().String(), err.Error())
					conn.Close()
				}
			}
		}
	}()

	srv.wgStartG.Wait()

	return nil
}

// CloseClient will close the client socket connection.
func (srv *Server) CloseClient(tcpAddr *net.TCPAddr) error {
	c, err := srv.clients.find(tcpAddr.IP.String())
	if err != nil {
		return fmt.Errorf("IP[ %s ] : disconnected", tcpAddr.String())
	}

	// Drop the connection using a goroutine since we are on the
	// socket goroutine most likely.
	go c.close()

	return nil
}

// Addr returns the listener's network address. This may be different than the values
// provided in the configuration, for example if configuration port value is 0.
func (srv *Server) Addr() net.Addr {
	return srv.tcpAddr
}

// Clients returns the number of active clients connected by user ID.
func (srv *Server) Clients() []string {
	clients := srv.clients.copy()

	users := make([]string, 0, len(clients))
	for _, c := range clients {
		users = append(users, c.UserID())
	}

	return users
}

// Groom drops connections that are not active for the specified duration.
func (srv *Server) Groom(d time.Duration) {
	client := srv.clients.copy()

	now := time.Now().UTC()
	for _, c := range client {
		sub := now.Sub(c.lastAct)
		if sub >= d {
			// TODO
			// This is a blocking call that waits for the socket goroutine
			// to report its done. This parallel call should work well since
			// there is no error handling needed.
			srv.log(srv.ctx, srv.name, EvtGroom, TypInfo, c.ipAddress, "Last[ %v ] Dur[ %v ]", c.lastAct.Format(time.RFC3339), sub)
			go c.close()
		}
	}
}

// =============================================================================

// startNewClient takes a new connection and adds it to the manager.
func (srv *Server) startNewClient(conn net.Conn) error {
	key := ipAddress(conn)

	c, err := newClient(key, srv.name, srv.log, srv.clients, srv.handlers, conn)
	if err != nil {
		return err
	}

	srv.clients.add(key, c)

	c.start()

	return nil
}
