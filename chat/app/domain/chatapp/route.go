package chatapp

import (
	"net/http"

	"github.com/ardanlabs/usdl/chat/business/chatbus"
	"github.com/ardanlabs/usdl/chat/foundation/logger"
	"github.com/ardanlabs/usdl/chat/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, log *logger.Logger, chatBus *chatbus.Business) {
	api := newApp(log, chatBus)

	app.HandlerFunc(http.MethodGet, "", "/connect", api.connect)
}
