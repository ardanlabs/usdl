// Package chatbus provides supports for chat activity.
package chatbus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ardanlabs/usdl/foundation/logger"
	"github.com/ardanlabs/usdl/foundation/signature"
	"github.com/ardanlabs/usdl/foundation/web"
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

	c1.Consume(c.listenBus(), jetstream.PullMaxMessages(1))

	const maxWait = 10 * time.Second
	c.ping(maxWait)

	return &c, nil
}

// Handshake performs the connection handshake protocol.
func (c *Business) Handshake(ctx context.Context, w http.ResponseWriter, r *http.Request) (User, error) {
	var ws websocket.Upgrader
	conn, err := ws.Upgrade(w, r, nil)
	if err != nil {
		return User{}, fmt.Errorf("upgrade: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte("HELLO")); err != nil {
		return User{}, fmt.Errorf("write message: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	usr := User{
		Conn:     conn,
		LastPing: time.Now(),
		LastPong: time.Now(),
	}

	msg, err := c.readMessage(ctx, usr)
	if err != nil {
		return User{}, fmt.Errorf("read message: %w", err)
	}

	if err := json.Unmarshal(msg, &usr); err != nil {
		return User{}, fmt.Errorf("unmarshal message: %w", err)
	}

	// -------------------------------------------------------------------------

	if err := c.storer.Add(ctx, usr); err != nil {
		defer conn.Close()
		if err := conn.WriteMessage(websocket.TextMessage, []byte("Already Connected")); err != nil {
			return User{}, fmt.Errorf("write message: %w", err)
		}
		return User{}, fmt.Errorf("add user: %w", err)
	}

	usr.Conn.SetPongHandler(c.pong(usr.ID))

	// -------------------------------------------------------------------------

	v := fmt.Sprintf("WELCOME %s", usr.Name)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(v)); err != nil {
		return User{}, fmt.Errorf("write message: %w", err)
	}

	c.log.Info(ctx, "chat-handshake", "status", "complete", "usr", usr)

	return usr, nil
}

// ListenClient waits for messages from users.
func (c *Business) ListenClient(ctx context.Context, from User) {
	for {
		msg, err := c.readMessage(ctx, from)
		if err != nil {
			if c.isCriticalError(ctx, err) {
				return
			}
			continue
		}

		var inMsg incomingMessage
		if err := json.Unmarshal(msg, &inMsg); err != nil {
			c.log.Info(ctx, "loc-unmarshal", "ERROR", err)
			continue
		}

		c.log.Info(ctx, "CLIENT: msg recv", "fromNonce", inMsg.FromNonce, "from", from.ID, "to", inMsg.ToID, "encrypted", inMsg.Encrypted, "message", inMsg.Msg)

		dataThatWasSign := struct {
			ToID      common.Address
			Msg       [][]byte
			FromNonce uint64
		}{
			ToID:      inMsg.ToID,
			Msg:       inMsg.Msg,
			FromNonce: inMsg.FromNonce,
		}

		id, err := signature.FromAddress(dataThatWasSign, inMsg.V, inMsg.R, inMsg.S)
		if err != nil {
			c.log.Info(ctx, "loc-fromAddress", "ERROR", err)
			continue
		}

		if id != from.ID.Hex() {
			c.log.Info(ctx, "loc-signature check", "status", "signature does not match")
			continue
		}

		to, err := c.storer.Retrieve(ctx, inMsg.ToID)
		if err != nil {
			switch {
			case errors.Is(err, ErrNotExists):
				c.log.Info(ctx, "loc-retrieve", "status", "user not found, sending over bus")
				if err := c.sendMessageBus(ctx, from, inMsg); err != nil {
					c.log.Info(ctx, "loc-bussend", "ERROR", err)
				}

			default:
				c.log.Info(ctx, "loc-retrieve", "ERROR", err)
			}

			continue
		}

		if err := c.sendMessageClient(from, to, inMsg.FromNonce, inMsg.Encrypted, inMsg.Msg); err != nil {
			c.log.Info(ctx, "loc-send", "ERROR", err)
		}

		c.log.Info(ctx, "LOC: msg sent over web socket", "from", from.ID, "to", inMsg.ToID)
	}
}

// =============================================================================

func (c *Business) listenBus() func(msg jetstream.Msg) {
	ctx := web.SetTraceID(context.Background(), uuid.New())

	f := func(msg jetstream.Msg) {
		defer msg.Ack()

		var busMsg busMessage
		if err := json.Unmarshal(msg.Data(), &busMsg); err != nil {
			c.log.Info(ctx, "bus-unmarshal", "ERROR", err)
			return
		}

		if busMsg.CapID == c.capID {
			return
		}

		c.log.Info(ctx, "BUS: msg recv", "fromNonce", busMsg.FromNonce, "from", busMsg.FromID, "to", busMsg.ToID, "encrypted", busMsg.Encrypted, "message", busMsg.Msg, "fromName", busMsg.FromName)

		dataThatWasSign := struct {
			ToID      common.Address
			Msg       [][]byte
			FromNonce uint64
		}{
			ToID:      busMsg.ToID,
			Msg:       busMsg.Msg,
			FromNonce: busMsg.FromNonce,
		}

		id, err := signature.FromAddress(dataThatWasSign, busMsg.V, busMsg.R, busMsg.S)
		if err != nil {
			c.log.Info(ctx, "bus-fromAddress", "ERROR", err)
			return
		}

		if id != busMsg.FromID.Hex() {
			c.log.Info(ctx, "bus-signature check", "status", "signature does not match")
			return
		}

		to, err := c.storer.Retrieve(ctx, busMsg.ToID)
		if err != nil {
			switch {
			case errors.Is(err, ErrNotExists):
				c.log.Info(ctx, "bus-retrieve", "status", "user not found")

			default:
				c.log.Info(ctx, "bus-retrieve", "ERROR", err)
			}

			return
		}

		from := User{
			ID:   busMsg.FromID,
			Name: busMsg.FromName,
		}

		if err := c.sendMessageClient(from, to, busMsg.FromNonce, busMsg.Encrypted, busMsg.Msg); err != nil {
			c.log.Info(ctx, "bus-send", "ERROR", err)
		}

		c.log.Info(ctx, "BUS: msg sent over web socket", "from", busMsg.FromID, "to", busMsg.ToID)
	}

	return f
}

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

func (c *Business) readMessage(ctx context.Context, usr User) ([]byte, error) {
	type response struct {
		msg []byte
		err error
	}

	ch := make(chan response, 1)

	go func() {
		_, msg, err := usr.Conn.ReadMessage()
		if err != nil {
			ch <- response{nil, err}
		}
		ch <- response{msg, nil}
	}()

	var resp response

	select {
	case <-ctx.Done():
		c.storer.Remove(ctx, usr.ID)
		usr.Conn.Close()
		return nil, ctx.Err()

	case resp = <-ch:
		if resp.err != nil {
			c.storer.Remove(ctx, usr.ID)
			usr.Conn.Close()
			return nil, resp.err
		}
	}

	return resp.msg, nil
}

func (c *Business) sendMessageClient(from User, to User, fromNonce uint64, encrypted bool, msg [][]byte) error {
	m := outgoingMessage{
		From: outgoingUser{
			ID:    from.ID,
			Name:  from.Name,
			Nonce: fromNonce,
		},
		Encrypted: encrypted,
		Msg:       msg,
	}

	if err := to.Conn.WriteJSON(m); err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	return nil
}

func (c *Business) sendMessageBus(ctx context.Context, from User, inMsg incomingMessage) error {
	busMsg := busMessage{
		CapID:           c.capID,
		FromID:          from.ID,
		FromName:        from.Name,
		incomingMessage: inMsg,
	}

	d, err := json.Marshal(busMsg)
	if err != nil {
		return fmt.Errorf("send marshal message: %w", err)
	}

	_, err = c.js.Publish(ctx, c.subject, d)
	if err != nil {
		return fmt.Errorf("send publish: %w", err)
	}

	return nil
}

func (c *Business) pong(id common.Address) func(appData string) error {
	f := func(appData string) error {
		ctx := web.SetTraceID(context.Background(), uuid.New())

		c.log.Debug(ctx, "*** PONG ***", "id", id, "status", "started")
		defer c.log.Debug(ctx, "*** PONG ***", "id", id, "status", "completed")

		usr, err := c.storer.UpdateLastPong(ctx, id)
		if err != nil {
			c.log.Info(ctx, "*** PONG ***", "id", id, "ERROR", err)
			return nil
		}

		sub := usr.LastPong.Sub(usr.LastPing)
		c.log.Debug(ctx, "*** PONG ***", "id", id, "status", "received", "sub", sub.String(), "ping", usr.LastPing.String(), "pong", usr.LastPong.String())

		return nil
	}

	return f
}

func (c *Business) ping(maxWait time.Duration) {
	ticker := time.NewTicker(maxWait)

	go func() {
		ctx := web.SetTraceID(context.Background(), uuid.New())

		for {
			<-ticker.C

			c.log.Debug(ctx, "*** PING ***", "status", "started")

			for id, conn := range c.storer.Connections() {
				sub := conn.LastPong.Sub(conn.LastPing)
				if sub > maxWait {
					c.log.Info(ctx, "*** PING ***", "ping", conn.LastPing.String(), "pong", conn.LastPong.Second(), "maxWait", maxWait, "sub", sub.String())
					c.storer.Remove(ctx, id)
					continue
				}

				c.log.Debug(ctx, "*** PING ***", "status", "sending", "id", id)

				if err := conn.Conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
					c.log.Info(ctx, "*** PING ***", "status", "failed", "id", id, "ERROR", err)
				}

				if err := c.storer.UpdateLastPing(ctx, id); err != nil {
					c.log.Info(ctx, "*** PING ***", "status", "failed", "id", id, "ERROR", err)
				}
			}

			c.log.Debug(ctx, "*** PING ***", "status", "completed")
		}
	}()
}
