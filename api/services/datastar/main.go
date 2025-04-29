package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ardanlabs/usdl/app/domain/datastarapp"
	"github.com/ardanlabs/usdl/foundation/client"
	"github.com/ardanlabs/usdl/foundation/client/storage/dbfile"
	"github.com/ardanlabs/usdl/foundation/logger"
	"github.com/ardanlabs/usdl/foundation/web"
)

const (
	url            = "ws://localhost:3000/connect"
	configFilePath = "zarf/client"
)

func main() {
	var log *logger.Logger

	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx).String()
	}

	log = logger.New(os.Stdout, logger.LevelInfo, "DATASTAR", traceIDFn)

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func run(_ context.Context, log *logger.Logger) error {
	id, err := client.NewID(configFilePath)
	if err != nil {
		return fmt.Errorf("id: %w", err)
	}

	db, err := dbfile.NewDB(configFilePath, id.MyAccountID)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	// -------------------------------------------------------------------------

	ds, err := datastarapp.NewApp(log, id.MyAccountID)
	if err != nil {
		return fmt.Errorf("web: %w", err)
	}

	// -------------------------------------------------------------------------

	app := client.NewApp(db, id, url, ds)
	defer app.Close()

	ds.SetApp(app)

	// -------------------------------------------------------------------------

	if err := app.Handshake(db.MyAccount()); err != nil {
		return fmt.Errorf("handshake: %w", err)
	}

	if err := app.Run(); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}
