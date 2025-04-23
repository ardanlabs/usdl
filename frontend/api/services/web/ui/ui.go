package ui

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ardanlabs/usdl/frontend/foundation/client"
	"github.com/benbjohnson/hashfs"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nats-io/nats.go"
	datastar "github.com/starfederation/datastar/sdk/go"
)

//go:embed static/*
var staticFS embed.FS

var staticSys = hashfs.NewFS(staticFS)

const webUpdateSubject = "web.update"

type WebUI struct {
	app              *client.App
	usernames        map[common.Address]string
	HasUnseenMessage map[common.Address]bool
	myAccountID      common.Address
	visibleUser      common.Address
	ns               *embeddednats.Server
	nc               *nats.Conn
	messages         map[common.Address][]client.Message
}

func New(ctx context.Context, myAccountID common.Address) (*WebUI, error) {
	ns, err := embeddednats.New(ctx, embeddednats.WithDirectory("./frontend/zarf/data/nats"))
	if err != nil {
		return nil, fmt.Errorf("starting NATS server: %w", err)
	}
	ns.WaitForServer()

	nc, err := ns.Client()
	if err != nil {
		return nil, fmt.Errorf("connecting to NATS server: %w", err)
	}

	ui := &WebUI{
		usernames:        map[common.Address]string{},
		myAccountID:      myAccountID,
		HasUnseenMessage: map[common.Address]bool{},
		ns:               ns,
		nc:               nc,
		messages:         map[common.Address][]client.Message{},
	}

	ui.loadContacts()

	return ui, nil
}

var zeroUser common.Address

func (ui *WebUI) currentMessages() []client.Message {
	if ui.visibleUser == zeroUser {
		return nil
	}
	u, err := ui.app.QueryContactByID(ui.visibleUser)
	if err != nil {
		return nil
	}
	return u.Messages
}

func (ui *WebUI) Run() error {
	portRaw := os.Getenv("PORT")
	if portRaw == "" {
		portRaw = "1337"
	}
	port, err := strconv.Atoi(portRaw)
	if err != nil {
		return fmt.Errorf("port: %w", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.FS(staticSys)).ServeHTTP(w, r)
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		PageChat(ui, ui.currentMessages()...).Render(ctx, w)
	})

	router.Post("/chat/changeUser/{userHexID}", func(w http.ResponseWriter, r *http.Request) {
		userHexID := chi.URLParam(r, "userHexID")
		if userHexID == "" {
			http.Error(w, "userHexID not found", http.StatusBadRequest)
			return
		}

		id := common.HexToAddress(userHexID)
		if id == zeroUser {
			http.Error(w, "userHexID not found", http.StatusBadRequest)
			return
		}

		ui.HasUnseenMessage[id] = false
		ui.visibleUser = id

		ui.nc.Publish(webUpdateSubject, []byte("update"))
	})

	router.Post("/chat/sendMessage", func(w http.ResponseWriter, r *http.Request) {
		type Store struct {
			Message string `json:"message"`
		}
		store := &Store{}
		if err := datastar.ReadSignals(r, store); err != nil {
			http.Error(w, fmt.Sprintf("reading signals: %s", err), http.StatusBadRequest)
			return
		}
		ui.app.SendMessageHandler(ui.visibleUser, []byte(store.Message))
	})

	router.Get("/chat/updates", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)

		ch := make(chan *nats.Msg, 1)
		defer close(ch)

		sub, err := ui.nc.ChanSubscribe(webUpdateSubject, ch)
		if err != nil {
			http.Error(w, fmt.Sprintf("subscribing to %s: %s", webUpdateSubject, err), http.StatusInternalServerError)
			return
		}
		defer sub.Unsubscribe()

		for {
			select {
			case <-r.Context().Done():
				return
			case <-ch:
				// It's possible a message came in that is not written to
				// storage. We need to deal with that.

				msgs := ui.currentMessages()
				sse.MergeFragmentTempl(ChatFragment(ui, msgs...))
			}
		}
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	log.Printf("Starting server on port   %d\n", port)
	return srv.ListenAndServe()
}

func (ui *WebUI) SetApp(app *client.App) {
	ui.app = app
	ui.loadContacts()

	contacts := app.Contacts()
	if len(contacts) > 0 {
		ui.visibleUser = contacts[0].ID
		log.Printf("Default contact: %s\n", ui.visibleUser.Hex())
	}
}

func (ui *WebUI) WriteText(msg client.Message) {
	if msg.ID != ui.visibleUser {
		ui.HasUnseenMessage[msg.ID] = true
	}

	// Not all messages are written to storage because it can be an error.
	// We need to deal with that. The current ui.messages are not really
	// being used at this time.

	ui.messages[msg.ID] = append(ui.messages[msg.ID], msg)
	ui.nc.Publish(webUpdateSubject, []byte("update"))
}

func (ui *WebUI) UpdateContact(id common.Address, name string) {
	ui.usernames[id] = name
}

func (ui *WebUI) loadContacts() {
	if ui.app == nil {
		return
	}
	for _, user := range ui.app.Contacts() {
		ui.usernames[user.ID] = user.Name
	}
}

func prettyPrintHex(id common.Address) string {
	return fmt.Sprintf("%s...", id.Hex()[0:6])
}

func prettyPrintUser(id common.Address, name string) string {
	return fmt.Sprintf("%s %s", prettyPrintHex(id), name)
}
