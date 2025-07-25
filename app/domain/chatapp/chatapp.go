// Package chatapp provides the application layer for the chat domain.
package chatapp

import (
	"context"
	"errors"
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

	a.chat.DropTCPConnection(ctx, usr.ID)

	return web.NewNoResponse()
}

func (a *app) state(ctx context.Context, r *http.Request) web.Encoder {
	connections := a.chat.TCPConnections(ctx)

	return stateResponse{
		TCPConnections: connections,
	}
}

func (a *app) tcpConnectDrop(ctx context.Context, r *http.Request) web.Encoder {
	var tcpConnReq tcpConnRequest
	if err := web.Decode(r, &tcpConnReq); err != nil {
		return errs.Newf(errs.InvalidArgument, "invalid request: %s", err)
	}

	tuiUserID := common.HexToAddress(tcpConnReq.TUIUserID)
	clientUserID := common.HexToAddress(tcpConnReq.ClientUserID)

	if err := a.chat.DialTCPConnection(ctx, tuiUserID, clientUserID, "tcp4", tcpConnReq.TCPHost); err != nil {
		if errors.Is(err, chatbus.ErrClientAlreadyConnected) {
			if err := a.chat.DropTCPConnection(ctx, tuiUserID); err != nil {
				return errs.Newf(errs.AlreadyExists, "failed to drop tcp connection: %s", err)
			}
			return tcpConnDropResponse{Connected: false, Message: "tcp connection dropped"}
		}

		return errs.Newf(errs.Internal, "failed to dial tcp connection: %s", err)
	}

	return tcpConnDropResponse{Connected: true, Message: "tcp connection established"}
}
