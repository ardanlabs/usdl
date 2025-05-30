package tcp

import (
	"fmt"
	"maps"
	"net"
	"sync"
)

type clients struct {
	log       internalLogger
	clients   map[string]*client
	clientsMu sync.RWMutex
}

func newClients(log internalLogger) *clients {
	return &clients{
		log:     log,
		clients: make(map[string]*client),
	}
}

func (clt *clients) count() int {
	clt.clientsMu.RLock()
	defer clt.clientsMu.RUnlock()

	return len(clt.clients)
}

func (clt *clients) copy() map[string]*client {
	clt.clientsMu.RLock()
	defer clt.clientsMu.RUnlock()

	clients := make(map[string]*client)
	maps.Copy(clients, clt.clients)

	return clients
}

func (clt *clients) add(client *client) {
	clt.clientsMu.Lock()
	defer clt.clientsMu.Unlock()

	clt.clients[client.ipAddress] = client
}

func (clt *clients) close(conn net.Conn) {
	clt.clientsMu.Lock()
	defer clt.clientsMu.Unlock()

	ipAddress := conn.RemoteAddr().String()

	if _, exists := clt.clients[ipAddress]; !exists {
		clt.log(EvtRemove, TypError, ipAddress, "already removed")
		return
	}

	delete(clt.clients, ipAddress)

	conn.Close()
}

func (clt *clients) find(tcpAddr *net.TCPAddr) (*client, error) {
	clt.clientsMu.RLock()
	defer clt.clientsMu.RUnlock()

	c, exists := clt.clients[tcpAddr.String()]
	if !exists {
		return nil, fmt.Errorf("IP[ %s ] : not found", tcpAddr.String())
	}

	return c, nil
}
