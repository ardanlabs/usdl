package tcp

import (
	"net"
	"sync"
)

type listener struct {
	listener   *net.TCPListener
	listenerMu sync.RWMutex
}

func newListener() *listener {
	return &listener{}
}

func (l *listener) tcpListener() *net.TCPListener {
	l.listenerMu.RLock()
	defer l.listenerMu.RUnlock()

	return l.listener
}

func (l *listener) reset() {
	l.listenerMu.Lock()
	defer l.listenerMu.Unlock()

	l.listener.Close()
	l.listener = nil
}

func (l *listener) start(network string, laddr *net.TCPAddr) (*net.TCPListener, error) {
	l.listenerMu.Lock()
	defer l.listenerMu.Unlock()

	listener, err := net.ListenTCP(network, laddr)
	if err != nil {
		return nil, err
	}

	l.listener = listener

	return listener, nil
}
