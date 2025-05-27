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
	"path/filepath"
	"sync"
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

	cfg := tcp.Config{
		NetType:     "tcp4",
		Addr:        ":0",
		ConnHandler: tcpConnHandler{},
		ReqHandler:  tcpReqHandler{},
		RespHandler: tcpRespHandler{},
		Logger:      f,
	}

	// Create a new TCP value.
	u, err := tcp.New("TEST", cfg)
	if err != nil {
		return fmt.Errorf("creating new TCP listener: %w", err)
	}

	if err := u.Start(); err != nil {
		return fmt.Errorf("starting the TCP listener: %w", err)
	}
	defer func() {
		fmt.Println("CALLING STOP")
		u.Stop()
	}()

	// -------------------------------------------------------------------------
	// CLIENT SIDE

	var conn net.Conn
	for i := range 10 {
		fmt.Println("Waiting for server to start...")
		time.Sleep(300 * time.Millisecond)

		fmt.Println("Try Client Conenction:", i+1)
		addr := u.Addr()
		if addr == nil {
			continue
		}
		conn, err = net.Dial("tcp4", addr.String())
		if err != nil {
			if i < 10 {
				continue
			}

			return fmt.Errorf("dialing a new TCP connection: %w", err)
		}
		defer conn.Close()
		break
	}

	bufReader := bufio.NewReader(conn)
	bufWriter := bufio.NewWriter(conn)

	fmt.Print("CLIENT: SEND:", "Hello\n")

	if _, err := bufWriter.WriteString("Hello\n"); err != nil {
		return fmt.Errorf("sending data to the connection: %w", err)
	}
	bufWriter.Flush()

	response, err := bufReader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading response from the connection: %w", err)
	}

	fmt.Println("CLIENT: RECV:", response)

	return nil
}

// =============================================================================

// tcpConnHandler is required to process data.
type tcpConnHandler struct{}

// Bind is called to init to reader and writer.
func (tch tcpConnHandler) Bind(conn net.Conn) (io.Reader, io.Writer) {
	return bufio.NewReader(conn), bufio.NewWriter(conn)
}

// tcpReqHandler is required to process client messages.
type tcpReqHandler struct{}

// Read implements the udp.ReqHandler interface. It is provided a request
// value to popular and a io.Reader that was created in the Bind above.
func (tcpReqHandler) Read(ipAddress string, reader io.Reader) ([]byte, int, error) {
	bufReader := reader.(*bufio.Reader)

	fmt.Println("SERVER: WAITING ON READ")

	// Read a small string to keep the code simple.
	line, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, 0, err
	}

	fmt.Println("SERVER: MESSAGE READ")

	return []byte(line), len(line), nil
}

// Process is used to handle the processing of the message.
func (tcpReqHandler) Process(r *tcp.Request) {
	fmt.Println("SERVER: CLIENT MESSAGE:", string(r.Data))

	resp := tcp.Response{
		TCPAddr: r.TCPAddr,
		Data:    []byte("GOT IT\n"),
		Length:  7,
	}

	fmt.Println("SERVER: SEND:", string(resp.Data))

	r.TCP.Send(r.Context, &resp)
}

type tcpRespHandler struct{}

// Write is provided the user-defined writer and the data to write.
func (tcpRespHandler) Write(r *tcp.Response, writer io.Writer) error {
	bufWriter := writer.(*bufio.Writer)
	if _, err := bufWriter.WriteString(string(r.Data)); err != nil {
		return err
	}

	return bufWriter.Flush()
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
