package datastarapp

import (
	"net/http"

	"github.com/ardanlabs/usdl/foundation/web"
)

func (a *App) routes(app *web.App) error {
	app.HandlerFunc(http.MethodPost, "", "/chat/sendMessage", a.sendMessage)
	app.HandlerFunc(http.MethodPost, "", "/chat/changeUser/{userHexID}", a.changeUser)
	app.HandlerFunc(http.MethodGet, "", "/chat/updates", a.chatUpdates)
	app.HandlerFunc(http.MethodGet, "", "/", a.renderPage)

	return nil
}
