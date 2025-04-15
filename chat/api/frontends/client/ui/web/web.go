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
	app         *app.App
	usernames   map[common.Address]string
	myAccountID common.Address

	visibleUser common.Address
}

func New(myAccountID common.Address) *WebUI {
	ui := &WebUI{
		usernames:   map[common.Address]string{},
		myAccountID: myAccountID,
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

		ui.visibleUser = id

	})

	router.Get("/chat/updates", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		t := time.NewTicker(100 * time.Millisecond)

		for {
			select {
			case <-r.Context().Done():
				return
			case <-t.C:
				log.Print("tick")
				msgs := ui.currentMessages()
				c := chatMessageFragment(msgs...)
				sse.MergeFragmentTempl(c)
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

	switch msg.ID {
	case common.Address{}:
		fmt.Fprintln(os.Stdout, "-----")
		fmt.Fprintf(os.Stdout, "%s: %s\n", msg.Name, string(msg.Content))

	case ui.app.ID():
		fmt.Fprintln(os.Stdout, "-----")
		fmt.Fprintf(os.Stdout, "%s: %s\n", msg.Name, string(msg.Content))

	default:
		// idx := ui.list.GetCurrentItem()

		// _, currentID := ui.list.GetItemText(idx)
		// if currentID == "" {
		// 	fmt.Fprintln(os.Stdout, "-----")
		// 	fmt.Fprintln(os.Stdout, "id not found: "+msg.ID.Hex())
		// 	return
		// }

		// if msg.ID.Hex() == currentID {
		// 	fmt.Fprintln(os.Stdout, "-----")
		// 	fmt.Fprintf(os.Stdout, "%s: %s\n", msg.Name, string(msg.Content))
		// 	return
		// }

		// for i := range ui.list.GetItemCount() {
		// 	name, idStr := ui.list.GetItemText(i)
		// 	if msg.ID.Hex() == idStr {
		// 		ui.list.SetItemText(i, "* "+name, idStr)
		// 		ui.tviewApp.Draw()
		// 		return
		// 	}
		// }
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
