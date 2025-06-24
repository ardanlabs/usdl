package chatbus

import (
	"encoding/json"
	"fmt"

	"github.com/ardanlabs/usdl/foundation/tcp"
)

func (b *Business) tcpSendMessage(clt *tcp.Client, from UIUser, inMsg uiIncomingMessage) error {
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

	if _, err := clt.Writer.Write(d); err != nil {
		return fmt.Errorf("send tcp publish: %w", err)
	}

	return nil
}
