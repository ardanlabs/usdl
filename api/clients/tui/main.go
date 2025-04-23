package main

import (
	"fmt"
	"os"

	"github.com/ardanlabs/usdl/api/clients/tui/ui"
	"github.com/ardanlabs/usdl/foundation/client"
	"github.com/ardanlabs/usdl/foundation/client/storage/dbfile"
)

const (
	url            = "ws://localhost:3000/connect"
	configFilePath = "zarf/client"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	id, err := client.NewID(configFilePath)
	if err != nil {
		return fmt.Errorf("id: %w", err)
	}

	db, err := dbfile.NewDB(configFilePath, id.MyAccountID)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	// -------------------------------------------------------------------------

	ui := ui.New(id.MyAccountID)

	// -------------------------------------------------------------------------

	app := client.NewApp(db, id, url, ui)
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
