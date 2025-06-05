package tcp

import (
	"context"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

// client represents a single networked connection.
type client struct {
	log       internalLogger
	clients   *clients
	handlers  Handlers
	conn      net.Conn
	tcpAddr   *net.TCPAddr
	ipAddress string
	isIPv6    bool
	reader    io.Reader
	writer    io.Writer
	wg        sync.WaitGroup
	timeConn  time.Time
	lastAct   time.Time
	nReads    int
	nWrites   int
}

// newClient creates a new client for an incoming connection.
func newClient(log internalLogger, clients *clients, handlers Handlers, conn net.Conn) *client {
	now := time.Now().UTC()

	// Ask the user to bind the reader and writer they want to
	// use for this connection.
	r, w := handlers.Bind(conn)

	// This will be a TCPAddr 100% of the time.
	raddr := conn.RemoteAddr().(*net.TCPAddr)

	c := client{
		log:       log,
		clients:   clients,
		handlers:  handlers,
		conn:      conn,
		tcpAddr:   raddr,
		ipAddress: raddr.IP.String() + ":" + strconv.Itoa(raddr.Port),
		isIPv6:    raddr.IP.To4() == nil,
		reader:    r,
		writer:    w,
		timeConn:  now,
		lastAct:   now,
	}

	return &c
}

func (c *client) start() {
	c.wg.Add(1)
	go c.read()
}

func (c *client) close() {
	c.conn.Close()
	c.wg.Wait()

	c.log(EvtDrop, TypInfo, c.ipAddress, "connection closed")
}

func (c *client) read() {
	c.log(EvtRead, TypInfo, c.ipAddress, "client G started")

	defer func() {
		c.log(EvtDrop, TypInfo, c.ipAddress, "client G disconnected")
		c.clients.close(c.conn)
		c.wg.Done()
	}()

close:
	for {
		// Wait for a message to arrive.
		data, length, err := c.handlers.Read(c.ipAddress, c.reader)
		c.lastAct = time.Now().UTC()
		c.nReads++

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
				IP:   c.tcpAddr.IP,
				Port: c.tcpAddr.Port,
				Zone: c.tcpAddr.Zone,
			},
			IsIPv6:  c.isIPv6,
			ReadAt:  c.lastAct,
			Context: context.Background(),
			Data:    data,
			Length:  length,
		}

		// Process the request on this goroutine that is
		// handling the socket connection.
		c.handlers.Process(&r, c.writer)
	}
}
