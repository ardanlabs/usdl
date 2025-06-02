package chatbus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ardanlabs/usdl/foundation/signature"
	"github.com/ardanlabs/usdl/foundation/web"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// UIHandshake performs the connection handshake protocol.
func (c *Business) UIHandshake(ctx context.Context, w http.ResponseWriter, r *http.Request) (User, error) {
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

	msg, err := c.uiReadMessage(ctx, usr)
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

	usr.Conn.SetPongHandler(c.uiPong(usr.ID))

	// -------------------------------------------------------------------------

	v := fmt.Sprintf("WELCOME %s", usr.Name)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(v)); err != nil {
		return User{}, fmt.Errorf("write message: %w", err)
	}

	c.log.Info(ctx, "chat-handshake", "status", "complete", "usr", usr)

	return usr, nil
}

// UIListen waits for messages from users.
func (c *Business) UIListen(ctx context.Context, from User) {
	for {
		msg, err := c.uiReadMessage(ctx, from)
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

		// BILL: We want the logic to match the order of precedence in the order
		// we think about sending a message: websocket, peer-to-peer, bus.

		// If the user is found, send the message directly to the user.

		to, err := c.storer.Retrieve(ctx, inMsg.ToID)
		if err == nil {
			c.log.Info(ctx, "LOC: msg sent over web socket", "from", from.ID, "to", inMsg.ToID)

			if err := c.uiSendMessage(from, to, inMsg.FromNonce, inMsg.Encrypted, inMsg.Msg); err != nil {
				c.log.Info(ctx, "loc-send", "ERROR", err)
			}

			continue
		}

		// TODO: IF WE HAVE A PEER TO PEER CONNECTION, WE SHOULD SEND THE MESSAGE

		// We don't have a web socket connection for the user, send the
		// message over the bus.

		if errors.Is(err, ErrNotExists) {
			c.log.Info(ctx, "loc-retrieve", "status", "user not found, sending over bus")

			if err := c.busSendMessage(ctx, from, inMsg); err != nil {
				c.log.Info(ctx, "loc-bussend", "ERROR", err)
			}

			continue
		}

		c.log.Info(ctx, "loc-retrieve", "ERROR", err)
	}
}

// =============================================================================

func (c *Business) uiReadMessage(ctx context.Context, usr User) ([]byte, error) {
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

func (c *Business) uiSendMessage(from User, to User, fromNonce uint64, encrypted bool, msg [][]byte) error {
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

func (c *Business) uiPing(maxWait time.Duration) {
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

func (c *Business) uiPong(id common.Address) func(appData string) error {
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
