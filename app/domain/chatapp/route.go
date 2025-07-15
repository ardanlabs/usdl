package chatapp

import (
	"net/http"

	"github.com/ardanlabs/usdl/business/domain/chatbus"
	"github.com/ardanlabs/usdl/foundation/logger"
	"github.com/ardanlabs/usdl/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, log *logger.Logger, chatBus *chatbus.Business, serverAddr string) {
	api := newApp(log, chatBus, serverAddr)

	app.HandlerFunc(http.MethodGet, "", "/connect", api.connect)
	app.HandlerFunc(http.MethodGet, "", "/state", api.state)
	app.HandlerFunc(http.MethodPost, "", "/tcpconnectdrop", api.tcpConnectDrop)
}
