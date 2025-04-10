package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ardanlabs/usdl/chat/api/frontends/client/app"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App interface {
	SendMessageHandler(to common.Address, msg []byte) error
	Contacts() []app.User
	QueryContactByID(id common.Address) (app.User, error)
}

type WebUI struct {
	srv *http.Server
}

func New(myAccountID common.Address) *WebUI {
	return nil
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

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Datastar!!"))
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	log.Printf("Starting server on port %d\n", port)
	return srv.ListenAndServe()
}

func (ui *WebUI) SetApp(app App) {

}

func (ui *WebUI) WriteText(id string, msg app.Message) {

}

func (ui *WebUI) UpdateContact(id string, name string) {

}
