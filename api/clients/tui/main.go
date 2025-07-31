package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ardanlabs/usdl/api/clients/tui/ui"
	"github.com/ardanlabs/usdl/api/clients/tui/ui/client"
	"github.com/ardanlabs/usdl/api/clients/tui/ui/client/storage/dbfile"
	"github.com/ardanlabs/usdl/app/sdk/auth"
	"github.com/ardanlabs/usdl/foundation/agents/ollamallm"
	"github.com/ardanlabs/usdl/foundation/keystore"
	"github.com/golang-jwt/jwt/v4"
)

const (
	url            = "localhost:3000"
	configFilePath = "zarf/client"
	keysFolder     = "zarf/client/id/"
	activeKID      = "key"
	issuer         = "usdl project"
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

	// -------------------------------------------------------------------------

	ks := keystore.New()

	if _, err := ks.LoadByFileSystem(os.DirFS(keysFolder)); err != nil {
		return fmt.Errorf("loading keys by fs: %w", err)
	}
	authCfg := auth.Config{
		Log:       nil,
		KeyLookup: ks,
		Issuer:    issuer,
	}

	ath, err := auth.New(authCfg)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	tkn, err := ath.GenerateToken(activeKID, auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id.MyAccountID.Hex(),
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	})
	if err != nil {
		return fmt.Errorf("generating token: %w", err)
	}

	fmt.Println("JWT:", tkn)

	// -------------------------------------------------------------------------

	db, err := dbfile.NewDB(configFilePath, id, tkn)
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

	app := client.NewApp(db, id, url, ui, tkn)
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
