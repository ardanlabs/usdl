package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ardanlabs/usdl/foundation/tcp"
)

const filePath = "backend/zarf/client"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// -------------------------------------------------------------------------
	// SERVER SIDE

	logger := func(ctx context.Context, name string, evt string, typ string, ipAddress string, format string, a ...any) {
		traceCtxID := tcp.GetTraceID(ctx)
		log.Printf("EVENT: %s, %s, %s, %s, [%s] %s", name, evt, typ, ipAddress, traceCtxID, fmt.Sprintf(format, a...))
	}

	cfg := tcp.ServerConfig{
		NetType:  "tcp4",
		Addr:     "0.0.0.0:3001",
		Handlers: tcpSrvHandlers{},
		Logger:   logger,
	}

	// Create a new TCP value.
	server, err := tcp.NewServer("SERVER", cfg)
	if err != nil {
		return fmt.Errorf("creating new TCP listener: %w", err)
	}

	serverErrors := make(chan error, 1)

	go func() {
		fmt.Println("***> START LISTENING")
		serverErrors <- server.Listen()
	}()

	// -------------------------------------------------------------------------
	// CLIENT SIDE

	go func() {
		fmt.Println("***> START CLIENT TEST")
		tcpClient(logger)
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		fmt.Println("shutdown started", "signal", sig)
		defer fmt.Println("shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

var cltWG sync.WaitGroup

func tcpClient(logger tcp.Logger) error {
	const netType = "tcp4"
	const addr = "0.0.0.0:3001"

	cfg := tcp.ClientConfig{
		Handlers: tcpCltHandlers{},
		Logger:   logger,
	}

	cm, err := tcp.NewClientManager("CLIENT-MANAGER", cfg)
	if err != nil {
		return fmt.Errorf("creating new TCP client manager: %w", err)
	}
	defer cm.Shutdown(context.Background())

	cltWG.Add(2)

	for i := range 2 {
		go func() {
			fmt.Println(i, "***> Waiting for server to start...")
			time.Sleep(300 * time.Millisecond)

			userID := fmt.Sprintf("user-%d", i)

			clt, err := cm.Dial(context.TODO(), userID, netType, addr)
			if err != nil {
				fmt.Println(i, "dialing a new TCP connection: %w", err)
				return
			}

			fmt.Print("***> CLIENT: SEND:", "Hello\n")

			if _, err := clt.Writer.Write([]byte("Hello\n")); err != nil {
				fmt.Println(i, "sending data to the connection: %w", err)
				return
			}
		}()
	}

	cltWG.Wait()

	return nil
}

// =============================================================================

type metadata struct {
	FileName string `json:"filename"`
	Size     int16  `json:"size"`
}

// =============================================================================

type tcpSrvHandlers struct {
	md       metadata
	buf      []byte
	isLength bool
	chunkLen int
	written  int
}

func NewTCPSrvHandlers() *tcpSrvHandlers {
	return &tcpSrvHandlers{
		isLength: true,
	}
}

func (h *tcpSrvHandlers) Bind(clt *tcp.Client) error {
	fmt.Println("***> SERVER: BIND", "TRACEDID", clt.TraceID(), clt.Conn.RemoteAddr().String(), clt.Conn.LocalAddr().String())

	clt.Conn.Write([]byte("Hello\n"))

	bufReader := bufio.NewReader(clt.Conn)
	line, _, err := bufReader.ReadLine()
	if err != nil {
		fmt.Println("reading line:", err)
		return err
	}

	var md metadata
	if err := json.Unmarshal(line, &md); err != nil {
		fmt.Println("unmarshalling metadata:", err)
		return err
	}

	h.md = md

	return nil
}

func (h *tcpSrvHandlers) Read(clt *tcp.Client) ([]byte, int, error) {
	fmt.Println("***> SERVER: WAITING ON READ", "TRACEDID", clt.TraceID())

	data := make([]byte, 4096)
	n, err := clt.Conn.Read(data)
	if err != nil {
		fmt.Println("reading line:", err)
		return nil, 0, err
	}

	h.buf = append(h.buf, data...)

	return data, n, nil
}

func (h *tcpSrvHandlers) Process(r *tcp.Request, clt *tcp.Client) {
	fmt.Println("***> SERVER: CLIENT MESSAGE:", "TRACEDID", clt.TraceID(), string(r.Data))

	for {
		// WE NEED TO KNOW IF WE EXPECT TO HAVE A LENGTH AT THE BEGINNING
		// OF THE BUFFER OR NOT.
		if h.isLength {
			uVal := binary.LittleEndian.Uint16(h.buf[:2])
			h.chunkLen = int(uVal)
			h.buf = h.buf[2:]
		}

		// WRITTEN 0 BYTES
		// CHUNKLEN 4096
		// BUF 1024

		// WRITTEN 1024 BYTES
		// CHUNKLEN 4096
		// BUF 512

		switch {
		case len(h.buf) < (h.chunkLen - h.written):
			fmt.Println("WRITE %d BYTES TO DISK 1", len(h.buf))
			h.buf = h.buf[len(h.buf):]
			h.written += len(h.buf)

		case len(h.buf) == (h.chunkLen - h.written):
			fmt.Println("WRITE %d BYTES TO DISK 2", len(h.buf))
			h.buf = h.buf[len(h.buf):]
			h.written = 0
			h.isLength = true

		case len(h.buf) > (h.chunkLen - h.written):
			diff := h.chunkLen - h.written
			fmt.Println("WRITE %d BYTES TO DISK 3", diff)
			h.buf = h.buf[diff:]
			h.written = 0
			h.isLength = true
		}
	}
}

func (tcpSrvHandlers) Drop(clt *tcp.Client) {
	fmt.Println("***> SERVER: CONNECTION CLOSED", "TRACEDID", clt.TraceID())
}

// =============================================================================

type tcpCltHandlers struct{}

func (tcpCltHandlers) Bind(clt *tcp.Client) error {
	fmt.Println("***> CLIENT: BIND", "TRACEDID", clt.TraceID(), clt.Conn.RemoteAddr().String(), clt.Conn.LocalAddr().String())
	clt.Reader = bufio.NewReader(clt.Conn)

	return nil
}

func (tcpCltHandlers) Read(clt *tcp.Client) ([]byte, int, error) {
	bufReader := clt.Reader.(*bufio.Reader)

	fmt.Println("***> CLIENT: WAITING ON READ", "TRACEDID", clt.TraceID())

	// Read a small string to keep the code simple.
	line, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, 0, err
	}

	fmt.Println("***> CLIENT: MESSAGE READ", "TRACEDID", clt.TraceID())

	return []byte(line), len(line), nil
}

func (tcpCltHandlers) Process(r *tcp.Request, clt *tcp.Client) {
	fmt.Println("***> CLIENT: SERVER MESSAGE:", "TRACEDID", clt.TraceID(), string(r.Data))

	cltWG.Done()
}

func (tcpCltHandlers) Drop(clt *tcp.Client) {
	fmt.Println("***> CLIENT: CONNECTION CLOSED", "TRACEDID", clt.TraceID())
}
