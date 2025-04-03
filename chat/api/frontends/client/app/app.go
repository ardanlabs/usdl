// Package app provides client app support.
package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ardanlabs/usdl/chat/foundation/signature"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
)

type MyAccount struct {
	ID   common.Address
	Name string
}

type User struct {
	ID           common.Address
	Name         string
	AppLastNonce uint64
	LastNonce    uint64
	Key          string
	Messages     [][]byte
}

type Storage interface {
	QueryContactByID(id common.Address) (User, error)
	InsertContact(id common.Address, name string) (User, error)
	InsertMessage(id common.Address, msg []byte) error
	UpdateAppNonce(id common.Address, nonce uint64) error
	UpdateContactNonce(id common.Address, nonce uint64) error
	UpdateContactKey(id common.Address, key string) error
}

type UI interface {
	Run() error
	WriteText(id string, msg []byte)
	UpdateContact(id string, name string)
}

// =============================================================================

type outgoingMessage struct {
	ToID      common.Address `json:"toID"`
	Encrypted bool           `json:"encrypted"`
	Msg       []byte         `json:"msg"`
	FromNonce uint64         `json:"fromNonce"`
	V         *big.Int       `json:"v"`
	R         *big.Int       `json:"r"`
	S         *big.Int       `json:"s"`
}

type usr struct {
	ID    common.Address `json:"id"`
	Name  string         `json:"name"`
	Nonce uint64         `json:"nonce"`
}

type incomingMessage struct {
	From      usr    `json:"from"`
	Encrypted bool   `json:"encrypted"`
	Msg       []byte `json:"msg"`
}

// =============================================================================

type App struct {
	db   Storage
	ui   UI
	id   ID
	url  string
	conn *websocket.Conn
}

func NewApp(db Storage, ui UI, id ID, url string) *App {
	return &App{
		db:  db,
		ui:  ui,
		id:  id,
		url: url,
	}
}

func (app *App) Close() error {
	if app.conn == nil {
		return nil
	}

	return app.conn.Close()
}

func (app *App) Run() error {
	return app.ui.Run()
}

func (app *App) Handshake(acct MyAccount) error {
	conn, _, err := websocket.DefaultDialer.Dial(app.url, nil)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	app.conn = conn

	// -------------------------------------------------------------------------

	_, msg, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	if string(msg) != "HELLO" {
		return fmt.Errorf("unexpected message: %s", msg)
	}

	// -------------------------------------------------------------------------

	user := struct {
		ID   common.Address
		Name string
	}{
		ID:   app.id.MyAccountID,
		Name: acct.Name,
	}

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	// -------------------------------------------------------------------------

	if _, _, err = conn.ReadMessage(); err != nil {
		return fmt.Errorf("read: %w", err)
	}

	// -------------------------------------------------------------------------

	go func() {
		app.ReceiveCapMessage(conn)
	}()

	return nil
}

// =============================================================================

func (app *App) ReceiveCapMessage(conn *websocket.Conn) {
	for {
		_, rawMsg, err := conn.ReadMessage()
		if err != nil {
			app.ui.WriteText("system", fmt.Appendf(nil, "read: %s", err))
			return
		}

		var inMsg incomingMessage
		if err := json.Unmarshal(rawMsg, &inMsg); err != nil {
			app.ui.WriteText("system", fmt.Appendf(nil, "unmarshal: %s", err))
			return
		}

		user, err := app.db.QueryContactByID(inMsg.From.ID)
		switch {
		case err != nil:
			user, err = app.db.InsertContact(inMsg.From.ID, inMsg.From.Name)
			if err != nil {
				app.ui.WriteText("system", fmt.Appendf(nil, "add contact: %s", err))
				return
			}

			app.ui.UpdateContact(inMsg.From.ID.Hex(), inMsg.From.Name)

		default:
			inMsg.From.Name = user.Name
		}

		// -----------------------------------------------------------------

		expNonce := user.LastNonce + 1
		if inMsg.From.Nonce != expNonce {
			app.ui.WriteText("system", fmt.Appendf(nil, "invalid nonce: possible security issue with contact: got: %d, exp: %d", inMsg.From.Nonce, expNonce))
			return
		}

		if err := app.db.UpdateContactNonce(inMsg.From.ID, expNonce); err != nil {
			app.ui.WriteText("system", fmt.Appendf(nil, "update app nonce: %s", err))
			return
		}

		// ---------------------------------------------------------------------

		decryptedMsg, encryptedMsg, err := app.preprocessRecvMessage(inMsg)
		if err != nil {
			app.ui.WriteText("system", fmt.Appendf(nil, "preprocess message: %s", err))
			return
		}

		// ---------------------------------------------------------------------

		if decryptedMsg[0] != '/' {
			decMsg := formatMessage("You", decryptedMsg)
			encMsg := formatMessage("You", encryptedMsg)

			if err := app.db.InsertMessage(inMsg.From.ID, encMsg); err != nil {
				app.ui.WriteText("system", fmt.Appendf(nil, "add message: %s", err))
				return
			}

			app.ui.WriteText(inMsg.From.ID.Hex(), decMsg)
		}
	}
}

func (app *App) SendMessageHandler(to common.Address, msg []byte) error {
	if app.conn == nil {
		return fmt.Errorf("no connection")
	}

	if len(msg) == 0 {
		return fmt.Errorf("message cannot be empty")
	}

	usr, err := app.db.QueryContactByID(to)
	if err != nil {
		return fmt.Errorf("query contact: %w", err)
	}

	// -------------------------------------------------------------------------

	nonce := usr.AppLastNonce + 1

	decryptedMsg, encryptedMsg, err := app.preprocessSendMessage(usr, msg)
	if err != nil {
		return fmt.Errorf("preprocess message: %w", err)
	}

	msg = decryptedMsg

	var encrypted bool
	if encryptedMsg != nil {
		encrypted = true
		msg = encryptedMsg
	}

	// -------------------------------------------------------------------------

	dataToSign := struct {
		ToID      common.Address
		Msg       []byte
		FromNonce uint64
	}{
		ToID:      to,
		Msg:       msg,
		FromNonce: nonce,
	}

	v, r, s, err := signature.Sign(dataToSign, app.id.PrivKeyECDSA)
	if err != nil {
		return fmt.Errorf("signing: %w", err)
	}

	outMsg := outgoingMessage{
		ToID:      to,
		Encrypted: encrypted,
		Msg:       msg,
		FromNonce: nonce,
		V:         v,
		R:         r,
		S:         s,
	}

	data, err := json.Marshal(outMsg)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := app.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	// -------------------------------------------------------------------------

	if err := app.db.UpdateAppNonce(to, nonce); err != nil {
		return fmt.Errorf("update app nonce: %w", err)
	}

	// -------------------------------------------------------------------------

	if decryptedMsg[0] != '/' {
		decMsg := formatMessage("You", decryptedMsg)
		encMsg := formatMessage("You", encryptedMsg)

		if err := app.db.InsertMessage(to, encMsg); err != nil {
			return fmt.Errorf("add message: %w", err)
		}

		app.ui.WriteText(to.Hex(), decMsg)
	}

	return nil
}

// =============================================================================

func (app *App) preprocessRecvMessage(inMsg incomingMessage) ([]byte, []byte, error) {
	msg := inMsg.Msg

	// -------------------------------------------------------------------------
	// Process Normal Message

	if msg[0] != '/' {
		if !inMsg.Encrypted {
			return msg, nil, nil
		}

		decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, app.id.PrivKeyRSA, []byte(msg))
		if err != nil {
			return nil, nil, fmt.Errorf("encrypting message: %w", err)
		}

		return decryptedData, inMsg.Msg, nil
	}

	// -------------------------------------------------------------------------
	// Process Commands

	msgStr := string(msg)

	parts := strings.Split(msgStr[1:], " ")
	if len(parts) < 2 {
		return nil, nil, fmt.Errorf("invalid command format: parts: %d", len(parts))
	}

	switch parts[0] {
	case "key":
		if err := app.db.UpdateContactKey(inMsg.From.ID, msgStr[5:]); err != nil {
			return nil, nil, fmt.Errorf("updating key: %w", err)
		}
		return []byte("** updated contact's key **"), nil, nil
	}

	return nil, nil, fmt.Errorf("unknown command")
}

func (app *App) preprocessSendMessage(usr User, msg []byte) ([]byte, []byte, error) {

	// -------------------------------------------------------------------------
	// Process Normal Message

	if msg[0] != '/' {
		if usr.Key == "" {
			return msg, nil, nil
		}

		publicKey, err := getPublicKey(usr.Key)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to read public key: %w", err)
		}

		encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg)
		if err != nil {
			return nil, nil, fmt.Errorf("encrypting message: %w", err)
		}

		return msg, encryptedData, nil
	}

	// -------------------------------------------------------------------------
	// Process Commands

	msgStr := string(msg)
	msgStr = strings.TrimSpace(msgStr)
	msgStr = strings.ToLower(msgStr)

	parts := strings.Split(msgStr[1:], " ")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid command format")
	}

	switch parts[0] {
	case "share":
		switch parts[1] {
		case "key":
			if app.id.PubKeyRSA == "" {
				return nil, nil, fmt.Errorf("no key to share")
			}

			return fmt.Appendf(nil, "/key %s", app.id.PubKeyRSA), nil, nil
		}
	}

	return nil, nil, fmt.Errorf("unknown command")
}

// =============================================================================

func getPublicKey(pemBlock string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemBlock))
	if block == nil {
		return nil, errors.New("invalid key: Key must be a PEM encoded PKCS1 or PKCS8 key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	return publicKey.(*rsa.PublicKey), nil
}
