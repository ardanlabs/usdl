package tcp

import "time"

// Stat represents a client statistic.
type Stat struct {
	IP       string
	Reads    int
	Writes   int
	TimeConn time.Time
	LastAct  time.Time
}

// ClientStats return details for all active clients.
func (t *Server) ClientStats() []Stat {
	clients := t.clients.copy()

	stats := make([]Stat, 0, len(clients))

	for _, c := range clients {
		stat := Stat{
			IP:       c.ipAddress,
			Reads:    c.nReads,
			Writes:   c.nWrites,
			TimeConn: c.timeConn,
			LastAct:  c.lastAct,
		}

		stats = append(stats, stat)
	}

	return stats
}
