// Package datastarapp provides the application layer for the datastar domain.
package datastarapp

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/usdl/app/sdk/mid"
	"github.com/ardanlabs/usdl/foundation/client"
	"github.com/ardanlabs/usdl/foundation/logger"
	"github.com/ardanlabs/usdl/foundation/web"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nats-io/nats.go"
)

const webUpdateSubject = "web.update"

//go:embed static
var staticFS embed.FS

type App struct {
	log              *logger.Logger
	clientApp        *client.App
	myAccountID      common.Address
	nc               *nats.Conn
	visibleUser      common.Address
	usernames        map[common.Address]string
	hasUnseenMessage map[common.Address]bool
	messages         map[common.Address][]client.Message
}

func NewApp(log *logger.Logger, myAccountID common.Address) (*App, error) {
	ctx := context.Background()

	log.Info(ctx, "startup", "status", "starting NATS server")

	ns, err := embeddednats.New(ctx, embeddednats.WithDirectory("zarf/data/nats"))
	if err != nil {
		return nil, fmt.Errorf("starting NATS server: %w", err)
	}
	ns.WaitForServer()

	log.Info(ctx, "startup", "status", "connect to NATS server")

	nc, err := ns.Client()
	if err != nil {
		return nil, fmt.Errorf("connecting to NATS server: %w", err)
	}

	log.Info(ctx, "startup", "status", "connected to NATS server")

	a := App{
		log:              log,
		myAccountID:      myAccountID,
		nc:               nc,
		usernames:        map[common.Address]string{},
		hasUnseenMessage: map[common.Address]bool{},
		messages:         map[common.Address][]client.Message{},
	}

	a.loadContacts()

	return &a, nil
}

func (a *App) Run() error {
	app := web.NewApp(
		a.log.Info,
		mid.Logger(a.log),
		mid.Errors(a.log),
		mid.Panics(),
	)

	if err := a.routes(app); err != nil {
		return fmt.Errorf("creating routes: %w", err)
	}

	app.FileServer(staticFS, "static", "/static/")

	api := http.Server{
		Addr:    ":1337",
		Handler: app,
	}

	// -------------------------------------------------------------------------

	ctx := context.Background()

	a.log.Info(ctx, "startup", "status", "initializing API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// -------------------------------------------------------------------------

	serverErrors := make(chan error, 1)

	go func() {
		a.log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		a.log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer a.log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func (a *App) SetApp(app *client.App) {
	a.clientApp = app
	a.loadContacts()

	contacts := app.Contacts()
	if len(contacts) > 0 {
		a.visibleUser = contacts[0].ID
		log.Printf("Default contact: %s\n", a.visibleUser.Hex())
	}
}

func (a *App) CurrentMessages() []client.Message {
	if a.visibleUser.Cmp(common.Address{}) == 0 {
		return nil
	}

	u, err := a.clientApp.QueryContactByID(a.visibleUser)
	if err != nil {
		return nil
	}

	return u.Messages
}

func (a *App) WriteText(msg client.Message) {
	if msg.ID != a.visibleUser {
		a.hasUnseenMessage[msg.ID] = true
	}

	// Not all messages are written to storage because it can be an error.
	// We need to deal with that. The current a.messages are not really
	// being used at this time.

	a.messages[msg.ID] = append(a.messages[msg.ID], msg)
	a.nc.Publish(webUpdateSubject, []byte("update"))
}

func (a *App) UpdateContact(id common.Address, name string) {
	a.usernames[id] = name
}

// =============================================================================

func (a *App) loadContacts() {
	if a.clientApp == nil {
		return
	}
	for _, user := range a.clientApp.Contacts() {
		a.usernames[user.ID] = user.Name
	}
}

func prettyPrintHex(id common.Address) string {
	return fmt.Sprintf("%s...", id.Hex()[0:6])
}

func prettyPrintUser(id common.Address, name string) string {
	return fmt.Sprintf("%s %s", prettyPrintHex(id), name)
}
