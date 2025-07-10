package tcp

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ctxKey int

const traceKey ctxKey = 1

func setTraceID(ctx context.Context, traceID uuid.UUID) context.Context {
	return context.WithValue(ctx, traceKey, traceID)
}

// GetTraceID retrieves the trace ID from the context.
func GetTraceID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(traceKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}
	}
	return v
}

// =============================================================================

// Client represents a single networked connection.
type Client struct {
	Conn      net.Conn
	Reader    io.Reader
	Writer    io.Writer
	ctx       context.Context
	key       string
	userID    string
	name      string
	log       internalLogger
	tcpAddr   *net.TCPAddr
	clients   *clients
	handlers  Handlers
	ipAddress string
	isIPv6    bool
	wg        sync.WaitGroup
	timeConn  time.Time
	lastAct   time.Time
	nReads    int
	nWrites   int
}

// Key returns the key for the client.
func (clt *Client) Key() string {
	return clt.key
}

// UserID returns the user ID for the client.
func (clt *Client) UserID() string {
	return clt.userID
}

// SetUserID sets the user ID for the client.
func (clt *Client) SetUserID(userID string) {
	clt.userID = userID
}

// Context returns the context for the client.
func (clt *Client) Context() context.Context {
	return clt.ctx
}

// TraceID returns the trace ID for the client.
func (clt *Client) TraceID() uuid.UUID {
	return GetTraceID(clt.ctx)
}

// newClient creates a new client for an incoming connection.
func newClient(key string, name string, log internalLogger, clients *clients, handlers Handlers, conn net.Conn) (*Client, error) {
	now := time.Now().UTC()

	// This will be a TCPAddr 100% of the time.
	raddr := conn.RemoteAddr().(*net.TCPAddr)

	clt := Client{
		Conn:      conn,
		Reader:    conn,
		Writer:    conn,
		ctx:       setTraceID(context.Background(), uuid.New()),
		key:       key,
		name:      name,
		log:       log,
		tcpAddr:   raddr,
		clients:   clients,
		handlers:  handlers,
		ipAddress: ipAddress(conn),
		isIPv6:    raddr.IP.To4() == nil,
		timeConn:  now,
		lastAct:   now,
	}

	// Inform the user we have a socket connection for a new client.
	if err := handlers.Bind(&clt); err != nil {
		return nil, err
	}

	return &clt, nil
}

// SetContext sets the context for the client.
func (clt *Client) SetContext(ctx context.Context) {
	clt.ctx = ctx
}

func (clt *Client) start() {
	clt.wg.Add(1)
	go clt.read()
}

func (clt *Client) close() {
	clt.Conn.Close()
	clt.wg.Wait()

	clt.log(clt.ctx, clt.name, EvtDrop, TypInfo, clt.ipAddress, "connection closed")
}

func (clt *Client) read() {
	clt.log(clt.ctx, clt.name, EvtRead, TypInfo, clt.ipAddress, "client G started")

	defer func() {
		clt.handlers.Drop(clt)

		if err := clt.clients.close(clt.key); err != nil {
			clt.log(clt.ctx, clt.name, EvtDrop, TypError, clt.ipAddress, "error closing client: %s", err)
		}

		clt.log(clt.ctx, clt.name, EvtDrop, TypInfo, clt.ipAddress, "client G disconnected")

		clt.wg.Done()
	}()

close:
	for {
		// Wait for a message to arrive.
		data, length, err := clt.handlers.Read(clt)
		clt.lastAct = time.Now().UTC()
		clt.nReads++

		if err != nil {
			// temporary is declared to test for the existence of
			// the method coming from the net package.
			type temporary interface {
				Temporary() bool
			}

			if e, ok := err.(temporary); ok {
				if !e.Temporary() {
					break close
				}
			}

			if err == io.EOF {
				break close
			}

			continue
		}

		// Create the request.
		r := Request{
			TCPAddr: &net.TCPAddr{
				IP:   clt.tcpAddr.IP,
				Port: clt.tcpAddr.Port,
				Zone: clt.tcpAddr.Zone,
			},
			IsIPv6:  clt.isIPv6,
			ReadAt:  clt.lastAct,
			Context: context.Background(),
			Data:    data,
			Length:  length,
		}

		// Process the request on this goroutine that is
		// handling the socket connection.
		clt.handlers.Process(&r, clt)
	}
}
