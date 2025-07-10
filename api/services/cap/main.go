package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/ardanlabs/usdl/app/sdk/mux"
	"github.com/ardanlabs/usdl/business/domain/chatbus"
	"github.com/ardanlabs/usdl/business/domain/chatbus/managers/uicltmgr"
	"github.com/ardanlabs/usdl/foundation/logger"
	"github.com/ardanlabs/usdl/foundation/tcp"
	"github.com/ardanlabs/usdl/foundation/web"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

/*
	CAP to CAP communication
		- Drop TCP connections with double click
		- Fix: Minor UI issues with arrows
		- Show in message which protocol is used

	Datafile transfer
		- Private stream
		- Encryption

	Terminate TLS Connections
		- TCP P2P

	Group Chat
		- Allow users to create groups

	Refactor client
		- Clear history button

	Write Tests
		- Unit tests
		- Integration tests
*/

var build = "develop"

func main() {
	var log *logger.Logger

	traceIDFn := func(ctx context.Context) string {
		traceID := web.GetTraceID(ctx)
		if traceID != uuid.Nil {
			return traceID.String()
		}

		traceID = tcp.GetTraceID(ctx)
		if traceID != uuid.Nil {
			return traceID.String()
		}

		return ""
	}

	log = logger.New(os.Stdout, logger.LevelInfo, "CAP", traceIDFn)

	// -------------------------------------------------------------------------

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			APIHost         string        `conf:"default:0.0.0.0:3000"`
		}
		NATS struct {
			Host       string `conf:"default:demo.nats.io"`
			Subject    string `conf:"default:ardanlabs-cap"`
			IDFilePath string `conf:"default:zarf/cap"`
		}
		TCP struct {
			ServerName string `conf:"default:tcp-server"`
			ClientName string `conf:"default:tcp-clientmanager"`
			NetType    string `conf:"default:tcp4"`
			Addr       string `conf:"default:0.0.0.0:4000"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "CAP",
		},
	}

	const prefix = "CAP"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	log.BuildInfo(ctx)

	// -------------------------------------------------------------------------
	// Cap ID

	capID, err := getCapID(cfg.NATS.IDFilePath)
	if err != nil {
		return fmt.Errorf("id file parse: %w", err)
	}

	log.Info(ctx, "startup", "status", "getting cap", "capID", capID)

	// -------------------------------------------------------------------------
	// UI Client Manager

	uiCltMgr := uicltmgr.New(log)

	// -------------------------------------------------------------------------
	// TCP Server

	tcpSrvLogger := func(ctx context.Context, name string, evt string, typ string, ipAddress string, format string, a ...any) {
		log.Info(ctx, "Tcp Server Event", "name", name, "evt", evt, "typ", typ, "ipAddress", ipAddress, "info", fmt.Sprintf(format, a...))
	}

	tcpSrvCfg := tcp.ServerConfig{
		NetType:  cfg.TCP.NetType,
		Addr:     cfg.TCP.Addr,
		Handlers: chatbus.NewServerHandlers(log, uiCltMgr),
		Logger:   tcpSrvLogger,
	}

	tcpSrv, err := tcp.NewServer(cfg.TCP.ServerName, tcpSrvCfg)
	if err != nil {
		return fmt.Errorf("tcp server: %w", err)
	}
	defer func() {
		log.Info(ctx, "TCP", "status", "starting TCP server shutdown")
		tcpSrv.Shutdown(ctx)
	}()

	tcpErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "TCP", "status", "starting TCP server", "addr", cfg.TCP.Addr)
		tcpErrors <- tcpSrv.Listen()
	}()

	// -------------------------------------------------------------------------
	// TCP Client Manager

	tcpCltLogger := func(ctx context.Context, name string, evt string, typ string, ipAddress string, format string, a ...any) {
		log.Info(ctx, "Tcp Client Event", "name", name, "evt", evt, "typ", typ, "ipAddress", ipAddress, "info", fmt.Sprintf(format, a...))
	}

	cfgCltCfg := tcp.ClientConfig{
		Handlers: chatbus.NewClientHandlers(log),
		Logger:   tcpCltLogger,
	}

	tcpCM, err := tcp.NewClientManager(cfg.TCP.ClientName, cfgCltCfg)
	if err != nil {
		return fmt.Errorf("tcp client manager: %w", err)
	}
	defer func() {
		log.Info(ctx, "TCP", "status", "starting TCP client manager shutdown")
		tcpCM.Shutdown(ctx)
	}()

	// -------------------------------------------------------------------------
	// NATS and ChatBus

	nc, err := nats.Connect(cfg.NATS.Host)
	if err != nil {
		return fmt.Errorf("nats connect: %w", err)
	}
	defer nc.Close()

	cfgBus := chatbus.Config{
		Log:         log,
		NATSConn:    nc,
		UICltMgr:    uiCltMgr,
		TCPCltMgr:   tcpCM,
		TCPServer:   tcpSrv,
		NATSSubject: cfg.NATS.Subject,
		CAPID:       capID,
	}

	chatBus, err := chatbus.NewBusiness(cfgBus)
	if err != nil {
		return fmt.Errorf("chat: %w", err)
	}

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfgMux := mux.Config{
		Log:        log,
		ChatBus:    chatBus,
		ServerAddr: cfg.TCP.Addr,
	}

	webAPI := mux.WebAPI(cfgMux)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      webAPI,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case err := <-tcpErrors:
		return fmt.Errorf("tcp server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop web server gracefully: %w", err)
		}

		if err := tcpSrv.Shutdown(ctx); err != nil {
			return fmt.Errorf("could not stop tcp server gracefully: %w", err)
		}
	}

	return nil
}

func getCapID(idFilePath string) (uuid.UUID, error) {
	fileName := filepath.Join(idFilePath, "cap.id")

	if _, err := os.Stat(fileName); err != nil {
		os.MkdirAll(idFilePath, os.ModePerm)

		f, err := os.Create(fileName)
		if err != nil {
			return uuid.UUID{}, fmt.Errorf("id file create: %w", err)
		}

		if _, err := f.WriteString(uuid.NewString()); err != nil {
			return uuid.UUID{}, fmt.Errorf("id file write: %w", err)
		}

		f.Close()
	}

	f, err := os.Open(fileName)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("id file open: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("id file read: %w", err)
	}

	capID, err := uuid.Parse(string(b))
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("id file parse: %w", err)
	}

	return capID, nil
}
