package tcp

import (
	"context"
	"io"
	"net"
	"strconv"
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
	traceID   string
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

// newClient creates a new client for an incoming connection.
func newClient(log internalLogger, clients *clients, handlers Handlers, conn net.Conn) *Client {
	now := time.Now().UTC()

	// This will be a TCPAddr 100% of the time.
	raddr := conn.RemoteAddr().(*net.TCPAddr)

	traceID := uuid.New()
	ctx := setTraceID(context.Background(), traceID)

	clt := Client{
		Conn:      conn,
		Reader:    conn,
		Writer:    conn,
		ctx:       ctx,
		traceID:   traceID.String(),
		log:       log,
		tcpAddr:   raddr,
		clients:   clients,
		handlers:  handlers,
		ipAddress: raddr.IP.String() + ":" + strconv.Itoa(raddr.Port),
		isIPv6:    raddr.IP.To4() == nil,
		timeConn:  now,
		lastAct:   now,
	}

	// Inform the user we have a socket connection for a
	// new client.
	handlers.Bind(ctx, &clt)

	return &clt
}

// SetContext sets the context for the client.
func (clt *Client) SetContext(ctx context.Context) {
	clt.ctx = ctx
}

// TraceID retrieves the trace ID for the client connection.
func (clt *Client) TraceID() string {
	return clt.traceID
}

func (clt *Client) start() {
	clt.wg.Add(1)
	go clt.read()
}

func (clt *Client) close() {
	clt.Conn.Close()
	clt.wg.Wait()

	clt.log(EvtDrop, TypInfo, clt.ipAddress, clt.traceID, "connection closed")
}

func (clt *Client) read() {
	clt.log(EvtRead, TypInfo, clt.ipAddress, clt.traceID, "client G started")

	defer func() {
		clt.log(EvtDrop, TypInfo, clt.ipAddress, clt.traceID, "client G disconnected")
		if err := clt.clients.close(clt.Conn); err != nil {
			clt.log(EvtDrop, TypError, clt.ipAddress, clt.traceID, "error closing client: %s", err)
		}
		clt.wg.Done()
	}()

close:
	for {
		// Wait for a message to arrive.
		data, length, err := clt.handlers.Read(clt.ctx, clt)
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
		clt.handlers.Process(clt.ctx, &r, clt)
	}
}
