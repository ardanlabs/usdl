// Package chat provides supports for chat activity.
package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ardanlabs/usdl/chat/app/sdk/errs"
	"github.com/ardanlabs/usdl/chat/foundation/logger"
	"github.com/ardanlabs/usdl/chat/foundation/web"
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

// Users defines the set of behavior for user management.
type Users interface {
	Add(ctx context.Context, usr User) error
	UpdateLastPing(ctx context.Context, userID uuid.UUID) error
	UpdateLastPong(ctx context.Context, userID uuid.UUID) (User, error)
	Remove(ctx context.Context, userID uuid.UUID)
	Connections() map[uuid.UUID]Connection
	Retrieve(ctx context.Context, userID uuid.UUID) (User, error)
}

// Chat represents a chat support.
type Chat struct {
	log      *logger.Logger
	js       jetstream.JetStream
	stream   jetstream.Stream
	consumer jetstream.Consumer
	id       string
	subject  string
	users    Users
}

// New creates a new chat support.
func New(log *logger.Logger, conn *nats.Conn, subject string, users Users) (*Chat, error) {
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, fmt.Errorf("nats create js: %w", err)
	}

	ctx := context.Background()

	// js.DeleteStream(ctx, subject)

	s1, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     subject,
		Subjects: []string{subject},
		MaxAge:   24 * time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("nats add js: %w", err)
	}

	id := "b243072d-453c-4825-865a-f5f2994de643" //uuid.NewString()

	c1, err := s1.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   id,
		AckPolicy: jetstream.AckExplicitPolicy,
		//DeliverPolicy: jetstream.DeliverNewPolicy,
	})
	if err != nil {
		return nil, fmt.Errorf("nats add consumer: %w", err)
	}

	c := Chat{
		log:      log,
		js:       js,
		stream:   s1,
		consumer: c1,
		id:       id,
		subject:  subject,
		users:    users,
	}

	c.listenBus()

	const maxWait = 10 * time.Second
	c.ping(maxWait)

	return &c, nil
}

// Shutdown cleans up the chat system.
func (c *Chat) Shutdown(ctx context.Context) error {
	//return c.stream.DeleteConsumer(ctx, c.id)
	return nil
}

// Handshake performs the connection handshake protocol.
func (c *Chat) Handshake(ctx context.Context, w http.ResponseWriter, r *http.Request) (User, error) {
	var ws websocket.Upgrader
	conn, err := ws.Upgrade(w, r, nil)
	if err != nil {
		return User{}, errs.Newf(errs.FailedPrecondition, "unable to upgrade to websocket")
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

	if err := c.users.Add(ctx, usr); err != nil {
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

// Listen waits for messages from users.
func (c *Chat) Listen(ctx context.Context, from User) {
	for {
		msg, err := c.readMessage(ctx, from)
		if err != nil {
			if c.isCriticalError(ctx, err) {
				return
			}
			continue
		}

		var inMsg inMessage
		if err := json.Unmarshal(msg, &inMsg); err != nil {
			c.log.Info(ctx, "loc-unmarshal", "ERROR", err)
			continue
		}

		c.log.Info(ctx, "LOC: msg recv", "from", from.ID, "to", inMsg.ToID, "message", inMsg.Msg)

		to, err := c.users.Retrieve(ctx, inMsg.ToID)
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

		if err := c.sendMessage(from, to, inMsg.Msg); err != nil {
			c.log.Info(ctx, "loc-send", "ERROR", err)
		}

		c.log.Info(ctx, "LOC: msg sent", "from", from.ID, "to", inMsg.ToID)
	}
}

// =============================================================================

func (c *Chat) isCriticalError(ctx context.Context, err error) bool {
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

		if errors.Is(err, jetstream.ErrConsumerDeleted) {
			c.log.Info(ctx, "chat-isCriticalError", "status", "nats consumer deleted")
			return true
		}

		c.log.Info(ctx, "chat-isCriticalError", "ERROR", err, "TYPE", fmt.Sprintf("%T", err))
		return false
	}
}

func (c *Chat) listenBus() {
	ctx := web.SetTraceID(context.Background(), uuid.New())

	go func() {
		for {
			msg, err := c.readMessageBus(ctx)
			if err != nil {
				if c.isCriticalError(ctx, err) {
					return
				}
				continue
			}

			var busMsg busMessage
			if err := json.Unmarshal(msg.Data(), &busMsg); err != nil {
				c.log.Info(ctx, "bus-unmarshal", "ERROR", err)
				continue
			}

			c.log.Info(ctx, "BUS: msg recv", "from", busMsg.FromID, "to", busMsg.ToID, "message", busMsg.Msg)

			to, err := c.users.Retrieve(ctx, busMsg.ToID)
			if err != nil {
				switch {
				case errors.Is(err, ErrNotExists):
					c.log.Info(ctx, "bus-retrieve", "status", "user not found")

				default:
					c.log.Info(ctx, "bus-retrieve", "ERROR", err)
				}

				continue
			}

			from := User{
				ID:   busMsg.FromID,
				Name: busMsg.FromName,
			}

			if err := c.sendMessage(from, to, busMsg.Msg); err != nil {
				c.log.Info(ctx, "bus-send", "ERROR", err)
			}

			c.log.Info(ctx, "BUS: msg sent", "from", busMsg.FromID, "to", busMsg.ToID)
		}
	}()
}

func (c *Chat) readMessage(ctx context.Context, usr User) ([]byte, error) {
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
		c.users.Remove(ctx, usr.ID)
		usr.Conn.Close()
		return nil, ctx.Err()

	case resp = <-ch:
		if resp.err != nil {
			c.users.Remove(ctx, usr.ID)
			usr.Conn.Close()
			return nil, resp.err
		}
	}

	return resp.msg, nil
}

func (c *Chat) readMessageBus(ctx context.Context) (jetstream.Msg, error) {
	type response struct {
		msg jetstream.Msg
		err error
	}

	ch := make(chan response, 1)

	go func() {
		msg, err := c.consumer.Next()
		if err != nil {
			ch <- response{nil, err}
		}
		ch <- response{msg, nil}
	}()

	var resp response

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	case resp = <-ch:
		if resp.err != nil {
			return nil, resp.err
		}
	}

	if err := resp.msg.Ack(); err != nil {
		return nil, fmt.Errorf("ack message: %w", err)
	}

	return resp.msg, nil
}

func (c *Chat) sendMessage(from User, to User, msg string) error {
	m := outMessage{
		From: User{
			ID:   from.ID,
			Name: from.Name,
		},
		Msg: msg,
	}

	if err := to.Conn.WriteJSON(m); err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	return nil
}

func (c *Chat) sendMessageBus(ctx context.Context, from User, inMsg inMessage) error {
	busMsg := busMessage{
		FromID:   from.ID,
		FromName: from.Name,
		ToID:     inMsg.ToID,
		Msg:      inMsg.Msg,
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

func (c *Chat) pong(id uuid.UUID) func(appData string) error {
	f := func(appData string) error {
		ctx := web.SetTraceID(context.Background(), uuid.New())

		c.log.Debug(ctx, "*** PONG ***", "id", id, "status", "started")
		defer c.log.Debug(ctx, "*** PONG ***", "id", id, "status", "completed")

		usr, err := c.users.UpdateLastPong(ctx, id)
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

func (c *Chat) ping(maxWait time.Duration) {
	ticker := time.NewTicker(maxWait)

	go func() {
		ctx := web.SetTraceID(context.Background(), uuid.New())

		for {
			<-ticker.C

			c.log.Debug(ctx, "*** PING ***", "status", "started")

			for id, conn := range c.users.Connections() {
				sub := conn.LastPong.Sub(conn.LastPing)
				if sub > maxWait {
					c.log.Info(ctx, "*** PING ***", "ping", conn.LastPing.String(), "pong", conn.LastPong.Second(), "maxWait", maxWait, "sub", sub.String())
					c.users.Remove(ctx, id)
					continue
				}

				c.log.Debug(ctx, "*** PING ***", "status", "sending", "id", id)

				if err := conn.Conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
					c.log.Info(ctx, "*** PING ***", "status", "failed", "id", id, "ERROR", err)
				}

				if err := c.users.UpdateLastPing(ctx, id); err != nil {
					c.log.Info(ctx, "*** PING ***", "status", "failed", "id", id, "ERROR", err)
				}
			}

			c.log.Debug(ctx, "*** PING ***", "status", "completed")
		}
	}()
}
