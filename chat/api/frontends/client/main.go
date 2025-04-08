package main

import (
	"fmt"
	"os"

	"github.com/ardanlabs/usdl/chat/api/frontends/client/app"
	"github.com/ardanlabs/usdl/chat/api/frontends/client/storage/dbfile"
	"github.com/ardanlabs/usdl/chat/api/frontends/client/ui/tui"
)

const (
	url            = "ws://localhost:3000/connect"
	configFilePath = "chat/zarf/client"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	id, err := app.NewID(configFilePath)
	if err != nil {
		return fmt.Errorf("id: %w", err)
	}

	db, err := dbfile.NewDB(configFilePath, id.MyAccountID)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	// -------------------------------------------------------------------------

	ui := tui.New(id.MyAccountID)

	// -------------------------------------------------------------------------

	app := app.NewApp(db, ui, id, url)
	defer app.Close()

	ui.SetApp(app)

	// -------------------------------------------------------------------------

	if err := app.Handshake(db.MyAccount()); err != nil {
		return fmt.Errorf("handshake: %w", err)
	}

	if err := app.Run(); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}
