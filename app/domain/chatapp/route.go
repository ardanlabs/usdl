package chatapp

import (
	"net/http"

	"github.com/ardanlabs/usdl/business/domain/chatbus"
	"github.com/ardanlabs/usdl/foundation/logger"
	"github.com/ardanlabs/usdl/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, log *logger.Logger, chatBus *chatbus.Business) {
	api := newApp(log, chatBus)

	app.HandlerFunc(http.MethodGet, "", "/connect", api.connect)
	app.HandlerFunc(http.MethodPost, "", "/p2p", api.p2p)
}
