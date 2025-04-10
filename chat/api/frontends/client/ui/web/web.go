package web

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ardanlabs/usdl/chat/api/frontends/client/app"
	"github.com/benbjohnson/hashfs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed static/*
var staticFS embed.FS

var staticSys = hashfs.NewFS(staticFS)

type WebUI struct {
	app         *app.App
	usernames   map[string]string
	myAccountID common.Address
	messages    []app.Message
}

func New(myAccountID common.Address) *WebUI {
	ui := &WebUI{
		usernames:   map[string]string{},
		myAccountID: myAccountID,
	}
	ui.loadContacts()
	return ui
}

func (ui *WebUI) Run() error {
	log.Printf("FOOO!!!!")

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
		log.Printf("foo")
		PageChat(ui).Render(ctx, w)
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
}

func (ui *WebUI) WriteText(msg app.Message) {
	log.Printf("WriteText: %s %s", msg.ID, msg)

	ui.messages = append(ui.messages, msg)

	// if _, ok := ui.usernames[id]; !ok {
	// 	ui.loadContacts()
	// }
}

func (ui *WebUI) UpdateContact(id string, name string) {
	ui.usernames[id] = name
}

func (ui *WebUI) loadContacts() {
	if ui.app == nil {
		return
	}
	for _, user := range ui.app.Contacts() {
		ui.usernames[user.ID.Hex()] = user.Name
	}
}

func prettyPrintHex(id common.Address) string {
	return fmt.Sprintf("%s...", id.Hex()[0:6])
}

func prettyPrintUser(id, name string) string {
	return fmt.Sprintf("%s %s", id, id)
}
