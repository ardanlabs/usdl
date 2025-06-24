// Package chatapp provides the application layer for the chat domain.
package chatapp

import (
	"context"
	"net/http"

	"github.com/ardanlabs/usdl/app/sdk/errs"
	"github.com/ardanlabs/usdl/business/domain/chatbus"
	"github.com/ardanlabs/usdl/foundation/logger"
	"github.com/ardanlabs/usdl/foundation/web"
)

type app struct {
	log  *logger.Logger
	chat *chatbus.Business
}

func newApp(log *logger.Logger, chat *chatbus.Business) *app {
	return &app{
		log:  log,
		chat: chat,
	}
}

func (a *app) connect(ctx context.Context, r *http.Request) web.Encoder {
	usr, err := a.chat.UIHandshake(ctx, web.GetWriter(ctx), r)
	if err != nil {
		return errs.Newf(errs.FailedPrecondition, "handshake failed: %s", err)
	}
	defer usr.UIConn.Close()

	a.chat.UIListen(ctx, usr)

	return web.NewNoResponse()
}

func (a *app) p2p(ctx context.Context, r *http.Request) web.Encoder {
	// WE NEED A DATA MODEL WITH USER INFORMATION
	// a.chat.DialTCPConnection(ctx, userID, "tcp4", "localhost:8080")

	return nil
}
