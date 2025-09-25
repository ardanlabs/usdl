package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/usdl/foundation/tcp"
)

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
		Handlers: NewTCPSrvHandlers(),
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

	tcpClient(logger)

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
	// defer cm.Shutdown(context.Background())

	fmt.Println("***> Waiting for server to start...")
	time.Sleep(300 * time.Millisecond)

	userID := "1234"

	clt, err := cm.Dial(context.TODO(), userID, netType, addr)
	if err != nil {
		fmt.Println("dialing a new TCP connection: %w", err)
		return err
	}

	data := []string{
		"I'm happy. ",
		"She exercises every morning. ",
		"His dog barks loudly. ",
		"My school starts at 8:00. ",
		"We always eat dinner together. ",
		"They take the bus to work. ",
		"He doesn't like vegetables. ",
		"I don't want anything to drink. ",
		"This little black dress isn't expensive. ",
		"Those kids don't speak English.",
	}

	var streamBuf bytes.Buffer
	for _, d := range data {
		l := len(d)
		lenD := make([]byte, 2)
		binary.LittleEndian.PutUint16(lenD, uint16(l))

		streamBuf.Write(lenD)
		streamBuf.WriteString(d)
	}

	stream := streamBuf.Bytes()

	fmt.Println("**> CLIENT SENDING BYTES")

	for len(stream) > 0 {
		v := min(len(stream), min(rand.Intn(20), len(stream)))
		if _, err := clt.Conn.Write([]byte(stream[:v])); err != nil {
			fmt.Println("**> CLIENT SENDING BYTES", err)
		}

		stream = stream[v:]
	}

	return nil
}

// =============================================================================

type metadata struct {
	FileName string `json:"filename"`
	Size     int16  `json:"size"`
}

// =============================================================================

type tcpSrvHandlers struct {
	md     metadata
	buffer []byte
	bufLen int
}

func NewTCPSrvHandlers() *tcpSrvHandlers {
	return &tcpSrvHandlers{
		buffer: make([]byte, 1024*32),
	}
}

func (h *tcpSrvHandlers) Bind(clt *tcp.Client) error {
	fmt.Println("***>1 SERVER: BIND", "TRACEDID", clt.TraceID(), clt.Conn.RemoteAddr().String(), clt.Conn.LocalAddr().String())

	clt.Conn.Write([]byte("Hello\n"))

	bufReader := bufio.NewReader(clt.Conn)
	line, _, err := bufReader.ReadLine()
	if err != nil {
		fmt.Println("reading line:", err)
		return err
	}

	fmt.Println("***>2 SERVER: BIND", "TRACEDID", clt.TraceID(), string(line))

	var md metadata
	if err := json.Unmarshal(line, &md); err != nil {
		fmt.Println("SERVER: unmarshalling metadata:", err)
		return err
	}

	h.md = md

	fmt.Println("***>3 SERVER: BIND", "TRACEDID", clt.TraceID(), clt.Conn.RemoteAddr().String(), clt.Conn.LocalAddr().String(), md)

	return nil
}

func (h *tcpSrvHandlers) Read(clt *tcp.Client) ([]byte, int, error) {
	fmt.Println("***> SERVER: WAITING ON READ", "TRACEDID", clt.TraceID())

	data := make([]byte, 4096)
	v, err := clt.Conn.Read(data)
	if err != nil {
		fmt.Printf("reading line: BYTES[%d] %v\n", v, err)
		return nil, 0, err
	}

	n := copy(h.buffer[h.bufLen:], data[:v])
	h.bufLen += n

	fmt.Println("***> SERVER: READ", "TRACEDID", clt.TraceID(), n)

	return data, n, nil
}

func (h *tcpSrvHandlers) Process(r *tcp.Request, clt *tcp.Client) {
	fmt.Println("***> SERVER: CLIENT MESSAGE:", "TRACEDID", clt.TraceID())

	// WE NEED A LOOP HERE FOR GETTING A BUNCH OF BYTES

	const prefix = 2

	// DO WE HAVE ENOUGH BYTES FOR LENGTH?
	if len(h.buffer) <= 1 {
		return
	}

	// DO WE HAVE ALL THE EXPECTED BYTES
	uVal := binary.LittleEndian.Uint16(h.buffer[:prefix])
	chunkLen := int(uVal)
	if chunkLen > h.bufLen-prefix {
		return
	}

	// WRITE BYTES TO DISK
	fmt.Println("***> WRITE BYTES TO DISK:", string(h.buffer[prefix:chunkLen+prefix]))

	// MOVE BYTES TO FRONT OF BUFFER
	n := copy(h.buffer, h.buffer[prefix+chunkLen:h.bufLen])
	h.bufLen = n
}

func (tcpSrvHandlers) Drop(clt *tcp.Client) {
	fmt.Println("***> SERVER: CONNECTION CLOSED", "TRACEDID", clt.TraceID())
}

// =============================================================================

type tcpCltHandlers struct{}

func (tcpCltHandlers) Bind(clt *tcp.Client) error {
	fmt.Println("***> CLIENT: BIND", "TRACEDID", clt.TraceID(), clt.Conn.RemoteAddr().String(), clt.Conn.LocalAddr().String())

	bufReader := bufio.NewReader(clt.Conn)
	line, _, err := bufReader.ReadLine()
	if err != nil {
		fmt.Println("reading line:", err)
		return err
	}

	fmt.Println("***> CLIENT: BIND", "TRACEDID", clt.TraceID(), clt.Conn.RemoteAddr().String(), clt.Conn.LocalAddr().String(), string(line))

	md := metadata{
		FileName: "bill.txt",
		Size:     10_000,
	}

	data, err := json.Marshal(md)
	if err != nil {
		fmt.Println("marshal:", err)
		return err
	}

	fmt.Println("***> CLIENT: BIND", "TRACEDID", clt.TraceID(), string(data))

	clt.Conn.Write(data)
	clt.Conn.Write([]byte("\n"))

	return nil
}

func (tcpCltHandlers) Read(clt *tcp.Client) ([]byte, int, error) {
	return nil, 0, nil
}

func (tcpCltHandlers) Process(r *tcp.Request, clt *tcp.Client) {
}

func (tcpCltHandlers) Drop(clt *tcp.Client) {
	fmt.Println("***> CLIENT: CONNECTION CLOSED", "TRACEDID", clt.TraceID())
}
