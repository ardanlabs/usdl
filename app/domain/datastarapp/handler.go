package datastarapp

import (
	"context"
	"net/http"

	"github.com/ardanlabs/usdl/app/sdk/errs"
	"github.com/ardanlabs/usdl/foundation/web"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nats-io/nats.go"
	datastar "github.com/starfederation/datastar/sdk/go"
)

func (a *App) renderPage(ctx context.Context, _ *http.Request) web.Encoder {
	w := web.GetWriter(ctx)
	PageChat(a, a.CurrentMessages()...).Render(ctx, w)

	return web.NewNoResponse()
}

func (a *App) changeUser(ctx context.Context, r *http.Request) web.Encoder {
	userHexID := web.Param(r, "userHexID")
	if userHexID == "" {
		return errs.Newf(errs.InvalidArgument, "userHexID not found")
	}

	a.log.Info(ctx, "changeUser", "userHexID", web.Param(r, "userHexID"))

	id := common.HexToAddress(userHexID)
	if id.Cmp(common.Address{}) == 0 {
		return errs.Newf(errs.InvalidArgument, "invalid userHexID")
	}

	a.hasUnseenMessage[id] = false
	a.visibleUser = id

	a.log.Info(ctx, "changeUser", "publish", "update")

	if err := a.nc.Publish(webUpdateSubject, []byte("update")); err != nil {
		return errs.New(errs.Internal, err)
	}

	return nil
}

func (a *App) sendMessage(_ context.Context, r *http.Request) web.Encoder {
	var store struct {
		Message string `json:"message"`
	}

	if err := datastar.ReadSignals(r, store); err != nil {
		return errs.New(errs.InternalOnlyLog, err)
	}

	a.clientApp.SendMessageHandler(a.visibleUser, []byte(store.Message))

	return nil
}

func (a *App) chatUpdates(ctx context.Context, r *http.Request) web.Encoder {
	w := web.GetWriter(ctx)

	sse := datastar.NewSSE(w, r)

	ch := make(chan *nats.Msg, 1)
	defer close(ch)

	a.log.Info(ctx, "chatUpdates", "subscribe", webUpdateSubject)

	sub, err := a.nc.ChanSubscribe(webUpdateSubject, ch)
	if err != nil {
		return errs.Newf(errs.Internal, "subscribing to %s: %s", webUpdateSubject, err)
	}
	defer sub.Unsubscribe()

	a.log.Info(ctx, "chatUpdates", "status", "SSE LOOP")

	for {
		select {
		case <-ctx.Done():
			a.log.Info(ctx, "chatUpdates", "status", "context done")
			return web.NewNoResponse()

		case v := <-ch:
			// It's possible a message came in that is not written to
			// storage. We need to deal with that.
			a.log.Info(ctx, "chatUpdates", "NATS", string(v.Data))

			msgs := a.CurrentMessages()
			sse.MergeFragmentTempl(ChatFragment(a, msgs...))
		}
	}
}
