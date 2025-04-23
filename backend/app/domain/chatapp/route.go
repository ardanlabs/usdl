package chatapp

import (
	"net/http"

	"github.com/ardanlabs/usdl/backend/business/chatbus"
	"github.com/ardanlabs/usdl/sdk/logger"
	"github.com/ardanlabs/usdl/sdk/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, log *logger.Logger, chatBus *chatbus.Business) {
	api := newApp(log, chatBus)

	app.HandlerFunc(http.MethodGet, "", "/connect", api.connect)
}
