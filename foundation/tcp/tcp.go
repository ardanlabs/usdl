package tcp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

// Set of event types.
const (
	EvtAccept = iota + 1
	EvtJoin
	EvtRead
	EvtRemove
	EvtDrop
	EvtGroom
	EvtStop
)

// Set of event sub types.
const (
	TypError = iota + 1
	TypInfo
)

var eventTypes = map[int]string{
	EvtAccept: "accept",
	EvtJoin:   "join",
	EvtRead:   "read",
	EvtRemove: "remove",
	EvtDrop:   "drop",
	EvtGroom:  "groom",
	EvtStop:   "stop",
}

var eventSubTypes = map[int]string{
	TypError: "error",
	TypInfo:  "info",
}

// =============================================================================

// Set of error variables for start up.
var (
	ErrInvalidNetType       = errors.New("invalid net type configuration")
	ErrInvalidHandlers      = errors.New("invalid handlers configuration")
	ErrInvalidLoggerHandler = errors.New("invalid logger handler configuration")
)

// =============================================================================

type internalLogger func(ctx context.Context, name string, evt int, typ int, ipAddress string, format string, a ...any)

// Logger defines an handler used to help log events.
type Logger func(ctx context.Context, name string, evt string, typ string, ipAddress string, format string, a ...any)

// =============================================================================

// Request is the message received by the client.
type Request struct {
	TCPAddr *net.TCPAddr
	IsIPv6  bool
	ReadAt  time.Time
	Context context.Context
	Data    []byte
	Length  int
}

// =============================================================================

// Errors provides support for multi client operations that might error.
type Errors []error

// Error implments the error interface for CltError.
func (ers Errors) Error() string {
	var b bytes.Buffer

	for _, err := range ers {
		b.WriteString(err.Error())
		b.WriteString("\n")
	}

	return b.String()
}

// =============================================================================

// Handlers is a collection of interfaces that must be implemented by the user.
type Handlers interface {
	ConnHandler
	ReqHandler
	DropHandler
}

// ConnHandler is implemented by the user to bind the connection
// to a reader and writer for processing.
type ConnHandler interface {

	// Bind is called to set the reader and writer.
	Bind(clt *Client) error
}

// ReqHandler is implemented by the user to implement the processing
// of request messages from the client.
type ReqHandler interface {

	// Read is provided a request and a user-defined reader for each client
	// connection on its own routine. Read must read a full request and return
	// the populated request value.
	// Returning io.EOF or a non temporary error will show down the connection.

	// Read is provided an ipaddress and the user-defined reader and must return
	// the data read off the wire and the length. Returning io.EOF or a non
	// temporary error will show down the listener.
	Read(clt *Client) ([]byte, int, error)

	// Process is used to handle the processing of the request.
	Process(r *Request, clt *Client)
}

type DropHandler interface {

	// Drop is called when a connection is dropped.
	Drop(clt *Client)
}

// =============================================================================

func ipAddress(conn net.Conn) string {
	tcpAddr := conn.RemoteAddr().(*net.TCPAddr)
	return fmt.Sprintf("%s:%d", tcpAddr.IP.String(), tcpAddr.Port)
}
