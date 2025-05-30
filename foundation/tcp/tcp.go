package tcp

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// TCP contains a set of networked client connections.
type TCP struct {
	name                   string
	log                    func(evt int, typ int, ipAddress string, format string, a ...any)
	netType                string
	addr                   string
	connHandler            ConnHandler
	reqHandler             ReqHandler
	respHandler            RespHandler
	ipAddress              string
	port                   int
	tcpAddr                *net.TCPAddr
	listener               *net.TCPListener
	listenerMu             sync.RWMutex
	clients                map[string]*client
	clientsMu              sync.RWMutex
	wgStartG               sync.WaitGroup
	shuttingDown           int32
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
		clients:     make(map[string]*client),
	}

	return &t, nil
}

// Name returns the name of the TCP manager.
func (t *TCP) Name() string {
	return t.name
}

// Start creates the accept routine and begins to accept connections.
func (t *TCP) Start() error {
	if t.tcpListener() != nil {
		return errors.New("this TCP has already been started")
	}

	t.wgStartG.Add(1)

	go func() {
		defer func() {
			t.log(EvtAccept, TypInfo, net.JoinHostPort(t.ipAddress, strconv.Itoa(t.port)), "shutdown")
			t.wgStartG.Done()
		}()

	startlistener:
		for {
			if atomic.LoadInt32(&t.shuttingDown) != 0 {
				t.resetTCPListener()
				break
			}

			listener, err := t.startTCPListener()
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
					if atomic.LoadInt32(&t.shuttingDown) != 0 {
						t.resetTCPListener()
						break startlistener
					}

					t.log(EvtAccept, TypError, conn.RemoteAddr().String(), err.Error())

					type temporary interface {
						Temporary() bool
					}

					if e, ok := err.(temporary); ok && !e.Temporary() {
						t.resetTCPListener()
						continue startlistener
					}

					continue
				}

				// Add this new connection to the manager map.
				t.join(conn)
			}
		}
	}()

	return nil
}

// Stop shuts down the manager and closes all connections.
func (t *TCP) Stop() error {
	t.log(EvtStop, TypInfo, "", "started")
	defer t.log(EvtStop, TypInfo, "", "completed")

	atomic.StoreInt32(&t.shuttingDown, 1)

	t.resetTCPListener()

	for _, c := range t.copyClients() {
		c.drop()
	}

	t.wgStartG.Wait()

	return nil
}

// Drop will close the socket connection.
func (t *TCP) Drop(tcpAddr *net.TCPAddr) error {

	// Find the client connection for this IPAddress.
	var c *client
	t.clientsMu.Lock()
	{
		// Validate this ipaddress and socket exists first.
		var ok bool
		if c, ok = t.clients[tcpAddr.String()]; !ok {
			t.clientsMu.Unlock()
			return fmt.Errorf("IP[ %s ] : disconnected", tcpAddr.String())
		}
	}
	t.clientsMu.Unlock()

	// Drop the connection using a goroutine since we are on the
	// socket goroutine most likely.
	go c.drop()
	return nil
}

// Send will deliver the response back to the client.
func (t *TCP) Send(ctx context.Context, r *Response) error {

	// Find the client connection for this IPAddress.
	var c *client
	t.clientsMu.Lock()
	{
		// Validate this ipaddress and socket exists first.
		var ok bool
		if c, ok = t.clients[r.TCPAddr.String()]; !ok {
			t.clientsMu.Unlock()
			return fmt.Errorf("IP[ %s ] : disconnected", r.TCPAddr.String())
		}

		// Increment the number of writes.
		c.nWrites++
	}
	t.clientsMu.Unlock()

	// Send the response.
	return t.respHandler.Write(r, c.writer)
}

// SendAll will deliver the response back to all connected clients.
func (t *TCP) SendAll(ctx context.Context, r *Response) error {
	var clts []*client
	t.clientsMu.Lock()
	{
		for _, c := range t.clients {
			clts = append(clts, c)
			c.nWrites++
		}
	}
	t.clientsMu.Unlock()

	// TODO: Consider doing this in parallel.
	var errors Errors
	for _, c := range clts {
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

// Connections returns the number of client connections.
func (t *TCP) Connections() int {
	var l int

	t.clientsMu.Lock()
	{
		l = len(t.clients)
	}
	t.clientsMu.Unlock()

	return l
}

// =============================================================================

func (t *TCP) tcpListener() *net.TCPListener {
	t.listenerMu.RLock()
	defer t.listenerMu.RUnlock()

	return t.listener
}

func (t *TCP) resetTCPListener() {
	t.listenerMu.Lock()
	defer t.listenerMu.Unlock()

	t.listener.Close()
	t.listener = nil
}

func (t *TCP) startTCPListener() (*net.TCPListener, error) {
	t.listenerMu.Lock()
	defer t.listenerMu.Unlock()

	listener, err := net.ListenTCP(t.netType, t.tcpAddr)
	if err != nil {
		return nil, err
	}

	t.listener = listener

	return listener, nil
}

func (t *TCP) copyClients() map[string]*client {
	t.clientsMu.RLock()
	defer t.clientsMu.RUnlock()

	clients := make(map[string]*client)
	maps.Copy(clients, t.clients)

	return clients
}

// =============================================================================

// Stat represents a client statistic.
type Stat struct {
	IP       string
	Reads    int
	Writes   int
	TimeConn time.Time
	LastAct  time.Time
}

// ClientStats return details for all active clients.
func (t *TCP) ClientStats() []Stat {
	var clts []*client
	t.clientsMu.Lock()
	{
		for _, v := range t.clients {
			clts = append(clts, v)
		}
	}
	t.clientsMu.Unlock()

	stats := make([]Stat, len(clts))
	for i, c := range clts {
		stats[i] = Stat{
			IP:       c.ipAddress,
			Reads:    c.nReads,
			Writes:   c.nWrites,
			TimeConn: c.timeConn,
			LastAct:  c.lastAct,
		}
	}

	return stats
}

// Clients returns the number of active clients connected.
func (t *TCP) Clients() int {
	var count int
	t.clientsMu.Lock()
	{
		count = len(t.clients)
	}
	t.clientsMu.Unlock()

	return count
}

// Groom drops connections that are not active for the specified duration.
func (t *TCP) Groom(d time.Duration) {
	var clts []*client
	t.clientsMu.Lock()
	{
		for _, v := range t.clients {
			clts = append(clts, v)
		}
	}
	t.clientsMu.Unlock()

	now := time.Now().UTC()
	for _, c := range clts {
		sub := now.Sub(c.lastAct)
		if sub >= d {

			// TODO
			// This is a blocking call that waits for the socket goroutine
			// to report its done. This parallel call should work well since
			// there is no error handling needed.
			t.log(EvtGroom, TypInfo, c.ipAddress, "Last[ %v ] Dur[ %v ]", c.lastAct.Format(time.RFC3339), sub)
			go c.drop()
		}
	}
}

// join takes a new connection and adds it to the manager.
func (t *TCP) join(conn net.Conn) {
	ipAddress := conn.RemoteAddr().String()
	t.log(EvtJoin, TypTrigger, ipAddress, "new connection")

	t.clientsMu.Lock()
	{
		// Validate this has not been joined already.
		if _, ok := t.clients[ipAddress]; ok {
			t.log(EvtJoin, TypError, ipAddress, "already connected")
			conn.Close()

			t.clientsMu.Unlock()
			return
		}

		// Add the client connection to the map.
		t.clients[ipAddress] = newClient(t, conn)
	}
	t.clientsMu.Unlock()
}

// remove deletes a connection from the manager.
func (t *TCP) remove(conn net.Conn) {
	ipAddress := conn.RemoteAddr().String()

	t.clientsMu.Lock()
	{
		// Validate this has not been removed already.
		if _, ok := t.clients[ipAddress]; !ok {
			t.log(EvtRemove, TypError, ipAddress, "already removed")
			t.clientsMu.Unlock()
			return
		}

		// Remove the client connection from the map.
		delete(t.clients, ipAddress)
	}
	t.clientsMu.Unlock()

	// Close the connection for safe keeping.
	conn.Close()
}
