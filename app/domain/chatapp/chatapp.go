// Package chatapp provides the application layer for the chat domain.
package chatapp

import (
	"context"
	"net/http"

	"github.com/ardanlabs/usdl/app/sdk/errs"
	"github.com/ardanlabs/usdl/business/domain/chatbus"
	"github.com/ardanlabs/usdl/foundation/logger"
	"github.com/ardanlabs/usdl/foundation/web"
	"github.com/ethereum/go-ethereum/common"
)

type app struct {
	log        *logger.Logger
	chat       *chatbus.Business
	serverAddr string
}

func newApp(log *logger.Logger, chat *chatbus.Business, serverAddr string) *app {
	return &app{
		log:        log,
		chat:       chat,
		serverAddr: serverAddr,
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

func (a *app) tcpConnect(ctx context.Context, r *http.Request) web.Encoder {
	var tcpConnReq tcpConnRequest
	if err := web.Decode(r, &tcpConnReq); err != nil {
		return errs.Newf(errs.InvalidArgument, "invalid request: %s", err)
	}

	userID := common.HexToAddress(tcpConnReq.UserID)

	if err := a.chat.DialTCPConnection(ctx, userID, "tcp4", tcpConnReq.TCPHost); err != nil {
		return errs.Newf(errs.InternalOnlyLog, "failed to dial tcp connection: %s", err)
	}

	return nil
}
