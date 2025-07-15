// Package client provides client app support.
package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ardanlabs/usdl/app/sdk/errs"
	"github.com/ardanlabs/usdl/foundation/signature"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
)

type MyAccount struct {
	ID          common.Address
	Name        string
	ProfilePath string
}

type Message struct {
	From        common.Address
	To          common.Address
	Name        string
	Content     [][]byte
	DateCreated time.Time
	Encrypted   bool
}

type User struct {
	ID           common.Address
	Name         string
	AppLastNonce uint64
	LastNonce    uint64
	Key          string
	TCPHost      string
	Messages     []Message
}

type Storage interface {
	Contacts() []User
	QueryContactByID(id common.Address) (User, error)
	InsertContact(id common.Address, name string) (User, error)
	InsertMessage(id common.Address, msg Message) error
	UpdateAppNonce(id common.Address, nonce uint64) error
	UpdateContactNonce(id common.Address, nonce uint64) error
	UpdateContactKey(id common.Address, key string) error
}

type UI interface {
	Run() error // Must be non-blocking
	WriteText(msg Message)
	AddContact(id common.Address, name string)
	ApplyContactPrefix(id common.Address, option string, add bool)
}

// =============================================================================

type outgoingMessage struct {
	ToID      common.Address `json:"toID"`
	Encrypted bool           `json:"encrypted"`
	Msg       [][]byte       `json:"msg"`
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
	From      usr      `json:"from"`
	Encrypted bool     `json:"encrypted"`
	Msg       [][]byte `json:"msg"`
}

// =============================================================================

type App struct {
	db   Storage
	ui   UI
	id   ID
	url  string
	conn *websocket.Conn
}

func NewApp(db Storage, id ID, url string, ui UI) *App {
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

func (app *App) ID() common.Address {
	return app.id.MyAccountID
}

func (app *App) Run() error {
	return app.ui.Run()
}

func (app *App) Handshake(acct MyAccount) error {
	url := fmt.Sprintf("ws://%s/connect", app.url)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
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
			app.ui.WriteText(errorMessage("read: %s", err))
			return
		}

		var inMsg incomingMessage
		if err := json.Unmarshal(rawMsg, &inMsg); err != nil {
			app.ui.WriteText(errorMessage("unmarshal: %s", err))
			return
		}

		user, err := app.db.QueryContactByID(inMsg.From.ID)
		switch {
		case err != nil:
			user, err = app.db.InsertContact(inMsg.From.ID, inMsg.From.Name)
			if err != nil {
				app.ui.WriteText(errorMessage("add contact: %s", err))
				return
			}

			app.ui.AddContact(inMsg.From.ID, inMsg.From.Name)

		default:
			inMsg.From.Name = user.Name
		}

		// -----------------------------------------------------------------

		if len(inMsg.Msg) == 0 {
			app.ui.WriteText(errorMessage("no message"))
			return
		}

		// -----------------------------------------------------------------

		if string(inMsg.Msg[0]) != "EVENT" {
			expNonce := user.LastNonce + 1
			if inMsg.From.Nonce < expNonce {
				app.ui.WriteText(errorMessage("invalid nonce: possible security issue with contact: got: %d, exp: %d", inMsg.From.Nonce, expNonce))
				return
			}

			if err := app.db.UpdateContactNonce(inMsg.From.ID, expNonce); err != nil {
				app.ui.WriteText(errorMessage("update app nonce: %s", err))
				return
			}
		}

		// ---------------------------------------------------------------------

		if err := app.preprocessRecvMessage(inMsg); err != nil {
			app.ui.WriteText(errorMessage("preprocess message: %s", err))
			return
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
	// Split message into chunks of 240 bytes so they can be encrypted.

	const maxBytes = 240
	var msgs [][]byte
	chunks := (len(msg) / maxBytes) + 1 // Calculate the number of chunks.
	var start int

	for range chunks {
		end := min(start+maxBytes, len(msg))
		msgs = append(msgs, msg[start:end])
		start = start + maxBytes
	}

	onWire, onScreen, err := app.preprocessSendMessage(usr, msgs)
	if err != nil {
		return fmt.Errorf("preprocess message: %w", err)
	}

	// -------------------------------------------------------------------------

	nonce := usr.AppLastNonce + 1

	dataToSign := struct {
		ToID      common.Address
		Msg       [][]byte
		FromNonce uint64
	}{
		ToID:      to,
		Msg:       onWire,
		FromNonce: nonce,
	}

	v, r, s, err := signature.Sign(dataToSign, app.id.PrivKeyECDSA)
	if err != nil {
		return fmt.Errorf("signing: %w", err)
	}

	var encrypted bool
	if usr.Key != "" {
		encrypted = true
	}

	outMsg := outgoingMessage{
		ToID:      to,
		Encrypted: encrypted,
		Msg:       onWire,
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

	if msg[0] != '/' {
		msg := Message{
			From:      app.id.MyAccountID,
			To:        to,
			Name:      "You",
			Content:   onScreen,
			Encrypted: encrypted,
		}

		if err := app.db.InsertMessage(to, msg); err != nil {
			return fmt.Errorf("add message: %w", err)
		}

		app.ui.WriteText(msg)
	}

	return nil
}

func (app *App) Contacts() []User {
	return app.db.Contacts()
}

func (app *App) QueryContactByID(id common.Address) (User, error) {
	return app.db.QueryContactByID(id)
}

// =============================================================================

// This provides a default client configuration, but it's recommended
// this is replaced by the user with application specific settings using
// the WithClient function at the time a AuthAPI is constructed.
// DualStack Deprecated: Fast Fallback is enabled by default. To disable, set FallbackDelay to a negative value.
var defaultClient = http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 15 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

type StateResponse struct {
	TCPConnections []common.Address `json:"tcp_connections"`
}

func (app *App) GetState(ctx context.Context) (StateResponse, error) {
	url := fmt.Sprintf("http://%s/state", app.url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return StateResponse{}, fmt.Errorf("create request error: %s: %w", url, err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")

	resp, err := defaultClient.Do(req)
	if err != nil {
		return StateResponse{}, fmt.Errorf("do: error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return StateResponse{}, fmt.Errorf("request error: status: %s: %w", http.StatusText(resp.StatusCode), err)
	}

	var state StateResponse
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return StateResponse{}, fmt.Errorf("decode: %w", err)
	}

	return state, nil
}

func (app *App) EstablishTCPConnection(ctx context.Context, tuiUserID common.Address, clientUserID common.Address) error {
	usr, err := app.db.QueryContactByID(clientUserID)
	if err != nil {
		return fmt.Errorf("query contact: %w", err)
	}

	if usr.TCPHost == "" {
		return fmt.Errorf("no TCP host found for contact: %s", clientUserID)
	}

	tcpConnRequest := struct {
		TUIUserID    string `json:"tui_user_id"`
		ClientUserID string `json:"client_user_id"`
		TCPHost      string `json:"tcp_host"`
	}{
		TUIUserID:    tuiUserID.String(),
		ClientUserID: clientUserID.String(),
		TCPHost:      usr.TCPHost,
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(tcpConnRequest); err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	url := fmt.Sprintf("http://%s/tcpconnect", app.url)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &b)
	if err != nil {
		return fmt.Errorf("create request error: %s: %w", url, err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")

	resp, err := defaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do: error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("copy error: %w", err)
	}

	var errs *errs.Error
	if err := json.Unmarshal(data, &errs); err != nil {
		return fmt.Errorf("failed: response: %s, decoding error: %w ", string(data), err)
	}

	return errs
}

// =============================================================================

func (app *App) preprocessRecvMessage(inMsg incomingMessage) error {
	msgs := inMsg.Msg

	if len(msgs) == 0 {
		return fmt.Errorf("no message")
	}

	// -------------------------------------------------------------------------
	// Process Event Message

	if string(msgs[0]) == "EVENT" {
		switch string(msgs[1]) {
		case "TCP-CONN":
			app.ui.WriteText(Message{
				From:    inMsg.From.ID,
				To:      app.id.MyAccountID,
				Name:    "system",
				Content: [][]byte{[]byte("TCP connection established from: " + inMsg.From.ID.String())},
			})

			app.ui.ApplyContactPrefix(inMsg.From.ID, "->", true)

		case "TCP-DROP":
			app.ui.WriteText(Message{
				From:    inMsg.From.ID,
				To:      app.id.MyAccountID,
				Name:    "system",
				Content: [][]byte{[]byte("TCP connection dropped from: " + inMsg.From.ID.String())},
			})

			app.ui.ApplyContactPrefix(inMsg.From.ID, "->", false)

		default:
			return fmt.Errorf("unknown event: %s", string(inMsg.Msg[0]))
		}

		return nil
	}

	// -------------------------------------------------------------------------
	// Process Normal Message

	if msgs[0][0] != '/' {
		if !inMsg.Encrypted {
			msg := Message{
				From:      inMsg.From.ID,
				To:        app.id.MyAccountID,
				Name:      inMsg.From.Name,
				Content:   msgs,
				Encrypted: false,
			}

			if err := app.db.InsertMessage(inMsg.From.ID, msg); err != nil {
				return fmt.Errorf("add message: %w", err)
			}

			app.ui.WriteText(msg)
			return nil
		}

		decryptedData := make([][]byte, len(msgs))

		for i, msg := range msgs {
			dd, err := rsa.DecryptPKCS1v15(rand.Reader, app.id.PrivKeyRSA, []byte(msg))
			if err != nil {
				return fmt.Errorf("decrypting message: %w", err)
			}

			decryptedData[i] = dd
		}

		msg := Message{
			From:      inMsg.From.ID,
			To:        app.id.MyAccountID,
			Name:      inMsg.From.Name,
			Content:   decryptedData,
			Encrypted: true,
		}

		if err := app.db.InsertMessage(inMsg.From.ID, msg); err != nil {
			return fmt.Errorf("add message: %w", err)
		}

		app.ui.WriteText(msg)
		return nil
	}

	// -------------------------------------------------------------------------
	// Process Commands

	msgStr := string(msgs[0])

	parts := strings.Split(msgStr[1:], " ")
	if len(parts) < 2 {
		return fmt.Errorf("invalid command format: parts: %d", len(parts))
	}

	switch parts[0] {
	case "key":
		if err := app.db.UpdateContactKey(inMsg.From.ID, msgStr[5:]); err != nil {
			return fmt.Errorf("updating key: %w", err)
		}

		msg := Message{
			From:      inMsg.From.ID,
			To:        app.id.MyAccountID,
			Name:      inMsg.From.Name,
			Content:   [][]byte{[]byte("** updated contact's key **")},
			Encrypted: false,
		}

		if err := app.db.InsertMessage(inMsg.From.ID, msg); err != nil {
			return fmt.Errorf("add message: %w", err)
		}

		app.ui.WriteText(msg)

		return nil
	}

	return fmt.Errorf("unknown command")
}

func (app *App) preprocessSendMessage(usr User, msgs [][]byte) (onWire [][]byte, onScreen [][]byte, err error) {

	// -------------------------------------------------------------------------
	// Process Normal Message

	if msgs[0][0] != '/' {
		if usr.Key == "" {
			return msgs, msgs, nil
		}

		publicKey, err := getPublicKey(usr.Key)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to read public key: %w", err)
		}

		encryptedData := make([][]byte, len(msgs))

		for i, msg := range msgs {
			ed, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg)
			if err != nil {
				return nil, nil, fmt.Errorf("encrypting message: %w", err)
			}

			encryptedData[i] = ed
		}

		return encryptedData, msgs, nil
	}

	// -------------------------------------------------------------------------
	// Process Commands

	msgStr := string(msgs[0])
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

			errMsg := fmt.Appendf(nil, "/key %s", app.id.PubKeyRSA)

			return [][]byte{errMsg}, [][]byte{errMsg}, nil
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
