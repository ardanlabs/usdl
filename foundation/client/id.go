package client

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const idFileName = "key.ecdsa"
const encFileName = "key.rsa"

type ID struct {
	MyAccountID  common.Address
	PrivKeyECDSA *ecdsa.PrivateKey
	PrivKeyRSA   *rsa.PrivateKey
	PubKeyRSA    string
}

func NewID(filePath string) (ID, error) {
	os.MkdirAll(filepath.Join(filePath, "id"), os.ModePerm)

	fileName := filepath.Join(filePath, "id", idFileName)

	var addr common.Address
	var pkECDSA *ecdsa.PrivateKey

	_, err := os.Stat(fileName)
	switch {
	case err != nil:
		addr, pkECDSA, err = createKeyID(fileName)

	default:
		addr, pkECDSA, err = readKeyID(fileName)
	}

	if err != nil {
		return ID{}, fmt.Errorf("id: %w", err)
	}

	// -------------------------------------------------------------------------

	fileName = filepath.Join(filePath, "id", encFileName)

	var pkRSA *rsa.PrivateKey

	_, err = os.Stat(fileName)
	switch {
	case err != nil:
		pkRSA, err = createKeyEnc(fileName)

	default:
		pkRSA, err = readKeyEnc(fileName)
	}

	if err != nil {
		return ID{}, fmt.Errorf("id: %w", err)
	}

	// -------------------------------------------------------------------------

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pkRSA.PublicKey)
	if err != nil {
		return ID{}, fmt.Errorf("marshaling public key: %w", err)
	}

	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, &publicBlock); err != nil {
		return ID{}, fmt.Errorf("encoding to public PEM: %w", err)
	}

	id := ID{
		MyAccountID:  addr,
		PrivKeyECDSA: pkECDSA,
		PrivKeyRSA:   pkRSA,
		PubKeyRSA:    buf.String(),
	}

	return id, nil
}

// =============================================================================

func createKeyID(fileName string) (common.Address, *ecdsa.PrivateKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("generateKey: %w", err)
	}

	if err := crypto.SaveECDSA(fileName, privateKey); err != nil {
		return common.Address{}, nil, fmt.Errorf("saveECDSA: %w", err)
	}

	return crypto.PubkeyToAddress(privateKey.PublicKey), privateKey, nil
}

func readKeyID(fileName string) (common.Address, *ecdsa.PrivateKey, error) {
	privateKey, err := crypto.LoadECDSA(fileName)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("loadECDSA: %w", err)
	}

	return crypto.PubkeyToAddress(privateKey.PublicKey), privateKey, nil
}

func createKeyEnc(fileName string) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generating key: %w", err)
	}

	privateFile, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("creating private file: %w", err)
	}
	defer privateFile.Close()

	privateBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return nil, fmt.Errorf("encoding to private file: %w", err)
	}

	return privateKey, nil
}

func readKeyEnc(fileName string) (*rsa.PrivateKey, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("opening key file: %w", err)
	}
	defer file.Close()

	pemData, err := io.ReadAll(io.LimitReader(file, 1024*1024))
	if err != nil {
		return nil, fmt.Errorf("reading auth private key: %w", err)
	}

	privatePEM := string(pemData)

	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		return nil, errors.New("invalid key: Key must be a PEM encoded PKCS1 or PKCS8 key")
	}

	var parsedKey any
	parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}

	pk, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("key is not a valid RSA private key")
	}

	return pk, nil
}
