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

func (b *Business) natsReadMessage() func(msg jetstream.Msg) {
	ctx := web.SetTraceID(context.Background(), uuid.New())

	f := func(msg jetstream.Msg) {
		defer msg.Ack()

		var natsMsg natsInOutMessage
		if err := json.Unmarshal(msg.Data(), &natsMsg); err != nil {
			b.log.Info(ctx, "nats-unmarshal", "ERROR", err)
			return
		}

		if natsMsg.CapID == b.capID {
			return
		}

		b.log.Info(ctx, "BUS: msg recv", "fromNonce", natsMsg.FromNonce, "from", natsMsg.FromID, "to", natsMsg.ToID, "encrypted", natsMsg.Encrypted, "message", natsMsg.Msg, "fromName", natsMsg.FromName)

		dataThatWasSign := struct {
			ToID      common.Address
			Msg       [][]byte
			FromNonce uint64
		}{
			ToID:      natsMsg.ToID,
			Msg:       natsMsg.Msg,
			FromNonce: natsMsg.FromNonce,
		}

		id, err := signature.FromAddress(dataThatWasSign, natsMsg.V, natsMsg.R, natsMsg.S)
		if err != nil {
			b.log.Info(ctx, "nats-fromAddress", "ERROR", err)
			return
		}

		if id != natsMsg.FromID.Hex() {
			b.log.Info(ctx, "nats-signature check", "status", "signature does not match")
			return
		}

		// If the user is found, send the message directly to the user.
		to, err := b.uiCltMgr.Retrieve(ctx, natsMsg.ToID)
		if err == nil {
			b.log.Info(ctx, "NATS: msg sent over web socket", "from", natsMsg.FromID, "to", natsMsg.ToID)

			from := UIUser{
				ID:   natsMsg.FromID,
				Name: natsMsg.FromName,
			}

			if err := uiSendMessage(from, to, natsMsg.FromNonce, natsMsg.Encrypted, natsMsg.Msg); err != nil {
				b.log.Info(ctx, "nats-send", "ERROR", err)
			}

			return
		}

		if !errors.Is(err, ErrNotExists) {
			b.log.Info(ctx, "nats-retrieve", "ERROR", err)
		}

		// We don't have a web socket connection for the user then drop the
		// message on the floor because we can't deliver it.

		b.log.Info(ctx, "nats-retrieve", "status", "user not found")
	}

	return f
}

func (b *Business) natsSendMessage(ctx context.Context, from UIUser, inMsg uiIncomingMessage) error {
	natsMsg := natsInOutMessage{
		CapID:             b.capID,
		FromID:            from.ID,
		FromName:          from.Name,
		uiIncomingMessage: inMsg,
	}

	d, err := json.Marshal(natsMsg)
	if err != nil {
		return fmt.Errorf("send nats marshal message: %w", err)
	}

	_, err = b.js.Publish(ctx, b.natsSubject, d)
	if err != nil {
		return fmt.Errorf("send nats publish: %w", err)
	}

	return nil
}
