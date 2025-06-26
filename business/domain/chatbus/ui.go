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
func (b *Business) UIHandshake(ctx context.Context, w http.ResponseWriter, r *http.Request) (UIUser, error) {
	var ws websocket.Upgrader
	conn, err := ws.Upgrade(w, r, nil)
	if err != nil {
		return UIUser{}, fmt.Errorf("upgrade: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte("HELLO")); err != nil {
		return UIUser{}, fmt.Errorf("write message: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	usr := UIUser{
		UIConn:   conn,
		LastPing: time.Now(),
		LastPong: time.Now(),
	}

	msg, err := b.uiReadMessage(ctx, usr)
	if err != nil {
		return UIUser{}, fmt.Errorf("read message: %w", err)
	}

	if err := json.Unmarshal(msg, &usr); err != nil {
		return UIUser{}, fmt.Errorf("unmarshal message: %w", err)
	}

	// Check that we have a valid user ID and Name.
	if usr.ID == (common.Address{}) || usr.Name == "" {
		defer conn.Close()
		if err := conn.WriteMessage(websocket.TextMessage, []byte("Invalid User ID or Name")); err != nil {
			return UIUser{}, fmt.Errorf("write message: %w", err)
		}
		return UIUser{}, fmt.Errorf("invalid user ID or name")
	}

	// -------------------------------------------------------------------------

	if err := b.uiCltMgr.Add(ctx, usr); err != nil {
		defer conn.Close()
		if err := conn.WriteMessage(websocket.TextMessage, []byte("Already Connected")); err != nil {
			return UIUser{}, fmt.Errorf("write message: %w", err)
		}
		return UIUser{}, fmt.Errorf("add user: %w", err)
	}

	usr.UIConn.SetPongHandler(b.uiPong(usr.ID))

	// -------------------------------------------------------------------------

	v := fmt.Sprintf("WELCOME %s", usr.Name)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(v)); err != nil {
		return UIUser{}, fmt.Errorf("write message: %w", err)
	}

	b.log.Info(ctx, "chat-handshake", "status", "complete", "usr", usr)

	return usr, nil
}

// UIListen waits for messages from users.
func (b *Business) UIListen(ctx context.Context, from UIUser) {
	for {
		msg, err := b.uiReadMessage(ctx, from)
		if err != nil {
			if b.isCriticalError(ctx, err) {
				return
			}
			continue
		}

		var inMsg uiIncomingMessage
		if err := json.Unmarshal(msg, &inMsg); err != nil {
			b.log.Info(ctx, "loc-unmarshal", "ERROR", err)
			continue
		}

		b.log.Info(ctx, "CLIENT: msg recv", "fromNonce", inMsg.FromNonce, "from", from.ID, "to", inMsg.ToID, "encrypted", inMsg.Encrypted, "message", inMsg.Msg)

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
			b.log.Info(ctx, "loc-fromAddress", "ERROR", err)
			continue
		}

		if id != from.ID.Hex() {
			b.log.Info(ctx, "loc-signature check", "status", "signature does not match")
			continue
		}

		// BILL: We want the logic to match the order of precedence in the order
		// we think about sending a message: websocket, peer-to-peer, nats.

		// ---------------------------------------------------------------------
		// Web Socket

		// If the user is found, send the message directly to the user.
		uiTo, err := b.uiCltMgr.Retrieve(ctx, inMsg.ToID)
		if err == nil {
			b.log.Info(ctx, "LOC: msg sent over web socket", "from", from.ID, "to", inMsg.ToID)

			if err := uiSendMessage(from, uiTo, inMsg.FromNonce, inMsg.Encrypted, inMsg.Msg); err != nil {
				b.log.Info(ctx, "loc-send", "ERROR", err)
			}

			continue
		}

		if !errors.Is(err, ErrNotExists) {
			b.log.Info(ctx, "loc-retrieve", "ERROR", err)
		}

		// ---------------------------------------------------------------------
		// TCP

		clt, err := b.tcpCltMgr.Retrieve(ctx, inMsg.ToID.String())
		if err == nil {
			b.log.Info(ctx, "LOC: msg sent over tcp", "from", from.ID, "to", inMsg.ToID)

			if err := b.tcpSendMessage(clt, from, inMsg); err != nil {
				b.log.Info(ctx, "loc-tcp-send", "ERROR", err)
			}
		}

		// ---------------------------------------------------------------------
		// NATS

		// We don't have a web socket connection for the user then send the
		// message over nats to find a CAP that can deliver the message.
		b.log.Info(ctx, "loc-retrieve", "status", "user not found, sending over nats")

		if err := b.natsSendMessage(ctx, from, inMsg); err != nil {
			b.log.Info(ctx, "loc-nats-send", "ERROR", err)
		}
	}
}

// =============================================================================

func (b *Business) uiReadMessage(ctx context.Context, usr UIUser) ([]byte, error) {
	type response struct {
		msg []byte
		err error
	}

	ch := make(chan response, 1)

	go func() {
		_, msg, err := usr.UIConn.ReadMessage()
		if err != nil {
			ch <- response{nil, err}
		}
		ch <- response{msg, nil}
	}()

	var resp response

	select {
	case <-ctx.Done():
		b.uiCltMgr.Remove(ctx, usr.ID)
		usr.UIConn.Close()
		return nil, ctx.Err()

	case resp = <-ch:
		if resp.err != nil {
			b.uiCltMgr.Remove(ctx, usr.ID)
			usr.UIConn.Close()
			return nil, resp.err
		}
	}

	return resp.msg, nil
}

func (b *Business) uiPing(maxWait time.Duration) {
	ticker := time.NewTicker(maxWait)

	go func() {
		ctx := web.SetTraceID(context.Background(), uuid.New())

		for {
			<-ticker.C

			b.log.Debug(ctx, "*** PING ***", "status", "started")

			for id, conn := range b.uiCltMgr.Connections() {
				sub := conn.LastPong.Sub(conn.LastPing)
				if sub > maxWait {
					b.log.Info(ctx, "*** PING ***", "ping", conn.LastPing.String(), "pong", conn.LastPong.Second(), "maxWait", maxWait, "sub", sub.String())
					b.uiCltMgr.Remove(ctx, id)
					continue
				}

				b.log.Debug(ctx, "*** PING ***", "status", "sending", "id", id)

				if err := conn.Conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
					b.log.Info(ctx, "*** PING ***", "status", "failed", "id", id, "ERROR", err)
				}

				if err := b.uiCltMgr.UpdateLastPing(ctx, id); err != nil {
					b.log.Info(ctx, "*** PING ***", "status", "failed", "id", id, "ERROR", err)
				}
			}

			b.log.Debug(ctx, "*** PING ***", "status", "completed")
		}
	}()
}

func (b *Business) uiPong(id common.Address) func(appData string) error {
	f := func(appData string) error {
		ctx := web.SetTraceID(context.Background(), uuid.New())

		b.log.Debug(ctx, "*** PONG ***", "id", id, "status", "started")
		defer b.log.Debug(ctx, "*** PONG ***", "id", id, "status", "completed")

		usr, err := b.uiCltMgr.UpdateLastPong(ctx, id)
		if err != nil {
			b.log.Info(ctx, "*** PONG ***", "id", id, "ERROR", err)
			return nil
		}

		sub := usr.LastPong.Sub(usr.LastPing)
		b.log.Debug(ctx, "*** PONG ***", "id", id, "status", "received", "sub", sub.String(), "ping", usr.LastPing.String(), "pong", usr.LastPong.String())

		return nil
	}

	return f
}

// =============================================================================

func uiSendMessage(from UIUser, to UIUser, fromNonce uint64, encrypted bool, msg [][]byte) error {
	m := uiOutgoingMessage{
		From: uiOutgoingUser{
			ID:    from.ID,
			Name:  from.Name,
			Nonce: fromNonce,
		},
		Encrypted: encrypted,
		Msg:       msg,
	}

	if err := to.UIConn.WriteJSON(m); err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	return nil
}
