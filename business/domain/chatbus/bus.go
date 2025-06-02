package chatbus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ardanlabs/usdl/foundation/signature"
	"github.com/ardanlabs/usdl/foundation/web"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

func (c *Business) busReadMessage() func(msg jetstream.Msg) {
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

		// BILL: We want the logic to match the order of precedence in the order
		// we think about sending a message: websocket, peer-to-peer, bus.

		// If the user is found, send the message directly to the user.

		to, err := c.storer.Retrieve(ctx, busMsg.ToID)
		if err == nil {
			c.log.Info(ctx, "BUS: msg sent over web socket", "from", busMsg.FromID, "to", busMsg.ToID)

			from := User{
				ID:   busMsg.FromID,
				Name: busMsg.FromName,
			}

			if err := c.uiSendMessage(from, to, busMsg.FromNonce, busMsg.Encrypted, busMsg.Msg); err != nil {
				c.log.Info(ctx, "bus-send", "ERROR", err)
			}

			return
		}

		// TODO: IF WE HAVE A PEER TO PEER CONNECTION, WE SHOULD SEND THE MESSAGE

		// We don't have a web socket connection for the user, ignore the message.

		if errors.Is(err, ErrNotExists) {
			c.log.Info(ctx, "bus-retrieve", "status", "user not found")

			return
		}

		c.log.Info(ctx, "bus-retrieve", "ERROR", err)
	}

	return f
}

func (c *Business) busSendMessage(ctx context.Context, from User, inMsg incomingMessage) error {
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
