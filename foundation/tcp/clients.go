package tcp

import (
	"errors"
	"fmt"
	"maps"
	"sync"
)

type clients struct {
	log       internalLogger
	clients   map[string]*Client
	clientsMu sync.RWMutex
}

func newClients(log internalLogger) *clients {
	return &clients{
		log:     log,
		clients: make(map[string]*Client),
	}
}

func (clt *clients) count() int {
	clt.clientsMu.RLock()
	defer clt.clientsMu.RUnlock()

	return len(clt.clients)
}

func (clt *clients) copy() map[string]*Client {
	clt.clientsMu.RLock()
	defer clt.clientsMu.RUnlock()

	clients := make(map[string]*Client)
	maps.Copy(clients, clt.clients)

	return clients
}

func (clt *clients) add(userID string, client *Client) {
	clt.clientsMu.Lock()
	defer clt.clientsMu.Unlock()

	clt.clients[userID] = client
}

func (clt *clients) close(userID string) error {
	clt.clientsMu.Lock()
	defer clt.clientsMu.Unlock()

	_, exists := clt.clients[userID]
	if !exists {
		return errors.New("already removed")
	}

	delete(clt.clients, userID)

	return nil
}

func (clt *clients) find(userID string) (*Client, error) {
	clt.clientsMu.RLock()
	defer clt.clientsMu.RUnlock()

	c, exists := clt.clients[userID]
	if !exists {
		return nil, fmt.Errorf("user[ %s ] : not found", userID)
	}

	return c, nil
}
