// Package chatbus provides supports for chat activity.
package chatbus

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/ardanlabs/usdl/foundation/logger"
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

// Storer defines the set of behavior for user management.
type Storer interface {
	Add(ctx context.Context, usr User) error
	UpdateLastPing(ctx context.Context, userID common.Address) error
	UpdateLastPong(ctx context.Context, userID common.Address) (User, error)
	Remove(ctx context.Context, userID common.Address)
	Connections() map[common.Address]Connection
	Retrieve(ctx context.Context, userID common.Address) (User, error)
}

// Business represents a chat support.
type Business struct {
	log      *logger.Logger
	js       jetstream.JetStream
	stream   jetstream.Stream
	consumer jetstream.Consumer
	capID    uuid.UUID
	subject  string
	storer   Storer
}

// NewBusiness creates a new chat support.
func NewBusiness(log *logger.Logger, conn *nats.Conn, users Storer, subject string, capID uuid.UUID) (*Business, error) {
	ctx := context.TODO()

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, fmt.Errorf("nats new js: %w", err)
	}

	// js.DeleteStream(ctx, subject)

	s1, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     subject,
		Subjects: []string{subject},
		MaxAge:   24 * time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("nats create js: %w", err)
	}

	c1, err := s1.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       capID.String(),
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
	})
	if err != nil {
		return nil, fmt.Errorf("nats create consumer: %w", err)
	}

	c := Business{
		log:      log,
		js:       js,
		stream:   s1,
		consumer: c1,
		capID:    capID,
		subject:  subject,
		storer:   users,
	}

	c1.Consume(c.busReadMessage(), jetstream.PullMaxMessages(1))

	const maxWait = 10 * time.Second
	c.uiPing(maxWait)

	return &c, nil
}

// =============================================================================

func (c *Business) isCriticalError(ctx context.Context, err error) bool {
	switch e := err.(type) {
	case *websocket.CloseError:
		c.log.Info(ctx, "chat-isCriticalError", "status", "client disconnected")
		return true

	case *net.OpError:
		if !e.Temporary() {
			c.log.Info(ctx, "chat-isCriticalError", "status", "client disconnected")
			return true
		}
		return false

	default:
		if errors.Is(err, context.Canceled) {
			c.log.Info(ctx, "chat-isCriticalError", "status", "client canceled")
			return true
		}

		if errors.Is(err, nats.ErrConnectionClosed) {
			c.log.Info(ctx, "chat-isCriticalError", "status", "nats connection closed")
			return true
		}

		c.log.Info(ctx, "chat-isCriticalError", "ERROR", err, "TYPE", fmt.Sprintf("%T", err))
		return false
	}
}
