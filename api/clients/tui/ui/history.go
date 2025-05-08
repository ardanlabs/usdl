package ui

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type data struct {
	msgs []string
	len  int
}

type history struct {
	cap     int
	history map[common.Address]data
	mu      sync.RWMutex
}

func NewHistory(cap int) *history {
	return &history{
		cap:     cap,
		history: make(map[common.Address]data),
	}
}

func (h *history) add(id common.Address, msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if d, ok := h.history[id]; !ok {
		d.msgs = make([]string, h.cap)
		h.history[id] = d
	}

	d := h.history[id]

	for i := range len(d.msgs) - 1 {
		d.msgs[i] = d.msgs[i+1]
	}

	d.msgs[len(d.msgs)-1] = msg

	if d.len < h.cap {
		d.len++
	}

	h.history[id] = d
}

func (h *history) retrieve(id common.Address) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	d := h.history[id]

	v := make([]string, d.len)

	from := len(d.msgs) - d.len
	to := len(d.msgs)

	copy(v, d.msgs[from:to])

	return v
}
