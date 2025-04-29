package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

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
