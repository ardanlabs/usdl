package web

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ardanlabs/usdl/chat/api/frontends/client/app"
	"github.com/benbjohnson/hashfs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	datastar "github.com/starfederation/datastar/sdk/go"
)

//go:embed static/*
var staticFS embed.FS

var staticSys = hashfs.NewFS(staticFS)

type WebUI struct {
	app              *app.App
	usernames        map[common.Address]string
	HasUnseenMessage map[common.Address]bool
	myAccountID      common.Address
	visibleUser      common.Address

	messages map[common.Address][]app.Message
}

func New(myAccountID common.Address) *WebUI {
	ui := &WebUI{
		usernames:        map[common.Address]string{},
		myAccountID:      myAccountID,
		HasUnseenMessage: map[common.Address]bool{},
		messages:         map[common.Address][]app.Message{},
	}
	ui.loadContacts()
	return ui
}

var zeroUser common.Address

func (ui *WebUI) currentMessages() []app.Message {
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
		t := time.NewTicker(100 * time.Millisecond)

		for {
			select {
			case <-r.Context().Done():
				return
			case <-t.C:
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

func (ui *WebUI) SetApp(app *app.App) {
	ui.app = app
	ui.loadContacts()

	contacts := app.Contacts()
	if len(contacts) > 0 {
		ui.visibleUser = contacts[0].ID
		log.Printf("Default contact: %s\n", ui.visibleUser.Hex())
	}
}

func (ui *WebUI) WriteText(msg app.Message) {
	if msg.ID != ui.visibleUser {
		ui.HasUnseenMessage[msg.ID] = true
	}

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
