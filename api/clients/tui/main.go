package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ardanlabs/usdl/api/clients/tui/ui"
	"github.com/ardanlabs/usdl/api/clients/tui/ui/client"
	"github.com/ardanlabs/usdl/api/clients/tui/ui/client/storage/dbfile"
	"github.com/ardanlabs/usdl/foundation/agents/ollamallm"
)

const (
	url            = "localhost:3000"
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

	db, err := dbfile.NewDB(configFilePath, id)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	// -------------------------------------------------------------------------

	agent, err := ollamallm.New(db.MyAccount().ProfilePath)
	if err != nil {
		return fmt.Errorf("ollama agent: %w", err)
	}

	// If we can't connect to the ollama agent, we can use it.
	fmt.Println("warming up the agent...")
	if _, err := agent.Chat(context.Background(), "warm up", nil); err != nil {
		agent = nil
	}

	// -------------------------------------------------------------------------

	ui := ui.New(id.MyAccountID, agent)

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
