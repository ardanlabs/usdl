// Package chatbus provides supports for chat activity.
package chatbus

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/ardanlabs/usdl/foundation/logger"
	"github.com/ardanlabs/usdl/foundation/tcp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Set of error variables.
var (
	ErrExists    = fmt.Errorf("user exists")
	ErrNotExists = fmt.Errorf("user doesn't exists")
)

// UIClientManager defines the set of behavior for user management.
type UIClientManager interface {
	Add(ctx context.Context, usr UIUser) error
	UpdateLastPing(ctx context.Context, userID common.Address) error
	UpdateLastPong(ctx context.Context, userID common.Address) (UIUser, error)
	Remove(ctx context.Context, userID common.Address)
	Connections() map[common.Address]UIConnection
	Retrieve(ctx context.Context, userID common.Address) (UIUser, error)
}

// TCPClientManager defines the set of behavior for user management.
type TCPClientManager interface {
	Dial(ctx context.Context, userID string, network string, address string) (*tcp.Client, error)
	Retrieve(ctx context.Context, userID string) (*tcp.Client, error)
}

type Config struct {
	Log         *logger.Logger
	NATSConn    *nats.Conn
	UICltMgr    UIClientManager
	TCPCltMgr   TCPClientManager
	NATSSubject string
	CAPID       uuid.UUID
}

// Business represents a chat support.
type Business struct {
	log         *logger.Logger
	js          jetstream.JetStream
	stream      jetstream.Stream
	consumer    jetstream.Consumer
	capID       uuid.UUID
	natsSubject string
	uiCltMgr    UIClientManager
	tcpCltMgr   TCPClientManager
}

// NewBusiness creates a new chat support.
func NewBusiness(cfg Config) (*Business, error) {
	ctx := context.TODO()

	js, err := jetstream.New(cfg.NATSConn)
	if err != nil {
		return nil, fmt.Errorf("nats new js: %w", err)
	}

	// js.DeleteStream(ctx, subject)

	s1, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     cfg.NATSSubject,
		Subjects: []string{cfg.NATSSubject},
		MaxAge:   24 * time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("nats create js: %w", err)
	}

	c1, err := s1.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       cfg.CAPID.String(),
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
	})
	if err != nil {
		return nil, fmt.Errorf("nats create consumer: %w", err)
	}

	b := Business{
		log:         cfg.Log,
		js:          js,
		stream:      s1,
		consumer:    c1,
		capID:       cfg.CAPID,
		natsSubject: cfg.NATSSubject,
		uiCltMgr:    cfg.UICltMgr,
		tcpCltMgr:   cfg.TCPCltMgr,
	}

	c1.Consume(b.natsReadMessage(), jetstream.PullMaxMessages(1))

	const maxWait = 10 * time.Second
	b.uiPing(maxWait)

	return &b, nil
}

// DialTCPConnection dials a tcp connection to the given address for a client
// tcp connection. The address should be in the format "host:port".
func (b *Business) DialTCPConnection(ctx context.Context, userID common.Address, network string, address string) error {
	_, err := b.tcpCltMgr.Dial(ctx, userID.String(), network, address)
	if err != nil {
		return fmt.Errorf("dial tcp connection: %w", err)
	}

	return nil
}

// =============================================================================

func (b *Business) isCriticalError(ctx context.Context, err error) bool {
	switch e := err.(type) {
	case *websocket.CloseError:
		b.log.Info(ctx, "chat-isCriticalError", "status", "client disconnected")
		return true

	case *net.OpError:
		if !e.Temporary() {
			b.log.Info(ctx, "chat-isCriticalError", "status", "client disconnected")
			return true
		}
		return false

	default:
		if errors.Is(err, context.Canceled) {
			b.log.Info(ctx, "chat-isCriticalError", "status", "client canceled")
			return true
		}

		if errors.Is(err, nats.ErrConnectionClosed) {
			b.log.Info(ctx, "chat-isCriticalError", "status", "nats connection closed")
			return true
		}

		b.log.Info(ctx, "chat-isCriticalError", "ERROR", err, "TYPE", fmt.Sprintf("%T", err))
		return false
	}
}
