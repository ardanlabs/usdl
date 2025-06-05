package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ardanlabs/usdl/foundation/tcp"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/nats-io/nats.go"
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

	f := func(evt string, typ string, ipAddress string, format string, a ...any) {
		log.Printf("EVENT: %s, %s, %s, %s", evt, typ, ipAddress, fmt.Sprintf(format, a...))
	}

	cfg := tcp.ServerConfig{
		NetType:  "tcp4",
		Addr:     "0.0.0.0:3001",
		Handlers: tcpHandlers{},
		Logger:   f,
	}

	// Create a new TCP value.
	server, err := tcp.NewServer("TEST", cfg)
	if err != nil {
		return fmt.Errorf("creating new TCP listener: %w", err)
	}

	serverErrors := make(chan error, 1)

	go func() {
		fmt.Println("***> START LISTENING")
		serverErrors <- server.Listen()
	}()

	go func() {
		fmt.Println("***> START CLIENT TEST")
		tcpClient(cfg)
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

func tcpClient(cfg tcp.ServerConfig) error {
	for i := range 2 {
		go func() {
			var conn net.Conn

			fmt.Println(i, "***> Waiting for server to start...")
			time.Sleep(300 * time.Millisecond)

			fmt.Println("***> Try Client Conenction:", i+1)
			if i == 9 {
				fmt.Println(i, "unable to connect to server after 10 attempts")
				return
			}

			var err error
			conn, err = net.Dial(cfg.NetType, cfg.Addr)
			if err != nil {
				fmt.Println(i, "dialing a new TCP connection: %w", err)
				return
			}
			defer conn.Close()

			if conn == nil {
				fmt.Println(i, "connection is nil, unable to connect to server")
				return
			}

			bufReader := bufio.NewReader(conn)
			bufWriter := bufio.NewWriter(conn)

			fmt.Print("***> CLIENT: SEND:", "Hello\n")

			if _, err := bufWriter.WriteString("Hello\n"); err != nil {
				fmt.Println(i, "sending data to the connection: %w", err)
				return
			}
			bufWriter.Flush()

			response, err := bufReader.ReadString('\n')
			if err != nil {
				fmt.Println(i, "reading response from the connection: %w", err)
				return
			}

			fmt.Println(i, "***> CLIENT: RECV:", response)
		}()
	}

	return nil
}

// =============================================================================

// tcpConnHandler is required to process data.
type tcpHandlers struct{}

// Bind is called to init to reader and writer.
func (tcpHandlers) Bind(clt *tcp.Client) {
	clt.Reader = bufio.NewReader(clt.Conn)
}

var bill atomic.Int64

// Read implements the udp.ReqHandler interface. It is provided a request
// value to popular and a io.Reader that was created in the Bind above.
func (tcpHandlers) Read(clt *tcp.Client) ([]byte, int, error) {
	bufReader := clt.Reader.(*bufio.Reader)

	fmt.Println("***> SERVER: WAITING ON READ")

	// Force a delay to simulate a long read.
	// if bill.CompareAndSwap(0, 1) {
	// 	time.Sleep(100 * time.Second)
	// }

	// Read a small string to keep the code simple.
	line, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, 0, err
	}

	fmt.Println("***> SERVER: MESSAGE READ")

	return []byte(line), len(line), nil
}

// Process is used to handle the processing of the message.
func (tcpHandlers) Process(r *tcp.Request, clt *tcp.Client) {
	fmt.Println("***> SERVER: CLIENT MESSAGE:", string(r.Data))

	resp := "GOT IT\n"

	fmt.Println("***> SERVER: SEND:", resp)

	if _, err := clt.Writer.Write([]byte(resp)); err != nil {
		fmt.Println("***> SERVER: ERROR SENDING RESPONSE:", err)
		return
	}
}

// =============================================================================

func natsexp() error {
	ctx := context.Background()

	ns, err := embeddednats.New(ctx, embeddednats.WithDirectory("zarf/data/nats"))
	if err != nil {
		return fmt.Errorf("starting NATS server: %w", err)
	}
	ns.WaitForServer()

	nc, err := ns.Client()
	if err != nil {
		return fmt.Errorf("connecting to NATS server: %w", err)
	}

	const webUpdateSubject = "web.update"

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		ch := make(chan *nats.Msg, 1)
		defer close(ch)

		sub, err := nc.ChanSubscribe(webUpdateSubject, ch)
		if err != nil {
			log.Printf("subscribing to %s: %s", webUpdateSubject, err)
			return
		}
		defer sub.Unsubscribe()

		log.Println("Waiting For message")

		v := <-ch
		log.Println(string(v.Data))
	}()

	log.Println("Publish message")

	time.Sleep(time.Second)

	if err := nc.Publish(webUpdateSubject, []byte("update")); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	wg.Wait()

	return nil
}

func key() error {
	fileName := filepath.Join(filePath, "key.rsa")

	if err := generatePrivateKey(fileName); err != nil {
		return fmt.Errorf("generatePrivateKey: %w", err)
	}

	// -------------------------------------------------------------------------

	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("opening key file: %w", err)
	}
	defer file.Close()

	pemData, err := io.ReadAll(io.LimitReader(file, 1024*1024))
	if err != nil {
		return fmt.Errorf("reading auth private key: %w", err)
	}

	privatePEM := string(pemData)

	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		return errors.New("invalid key: Key must be a PEM encoded PKCS1 or PKCS8 key")
	}

	var parsedKey any
	parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return err
		}
	}

	pk, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return errors.New("key is not a valid RSA private key")
	}

	// -------------------------------------------------------------------------

	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, &pk.PublicKey, []byte("Hi Kevin, this is a secret message!"))
	if err != nil {
		return fmt.Errorf("encrypting message: %w", err)
	}

	fmt.Println(string(encryptedData))
	fmt.Println("")

	// -------------------------------------------------------------------------

	decryptedData, err := rsa.DecryptPKCS1v15(nil, pk, encryptedData)
	if err != nil {
		return fmt.Errorf("decrypting message: %w", err)
	}

	fmt.Println(string(decryptedData))

	return nil
}

func generatePrivateKey(fileName string) error {
	if _, err := os.Stat(fileName); err == nil {
		return nil
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generating key: %w", err)
	}

	privateFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("creating private file: %w", err)
	}
	defer privateFile.Close()

	privateBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return fmt.Errorf("encoding to private file: %w", err)
	}

	return nil
}
