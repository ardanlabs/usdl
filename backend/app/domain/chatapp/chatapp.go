// Package chatapp provides the application layer for the chat service.
package chatapp

import (
	"context"
	"net/http"

	"github.com/ardanlabs/usdl/backend/app/sdk/errs"
	"github.com/ardanlabs/usdl/backend/business/chatbus"
	"github.com/ardanlabs/usdl/sdk/logger"
	"github.com/ardanlabs/usdl/sdk/web"
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
	usr, err := a.chat.Handshake(ctx, web.GetWriter(ctx), r)
	if err != nil {
		return errs.Newf(errs.FailedPrecondition, "handshake failed: %s", err)
	}
	defer usr.Conn.Close()

	a.chat.ListenClient(ctx, usr)

	return web.NewNoResponse()
}
