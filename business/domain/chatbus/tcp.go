package chatbus

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ardanlabs/usdl/foundation/logger"
	"github.com/ardanlabs/usdl/foundation/signature"
	"github.com/ardanlabs/usdl/foundation/tcp"
	"github.com/ethereum/go-ethereum/common"
)

// ClientHandlers implements the Handlers interface for the TCP client manager.
type ClientHandlers struct {
	log *logger.Logger
}

// NewClientHandlers creates a new instance of ClientHandlers.
func NewClientHandlers(log *logger.Logger) *ClientHandlers {
	return &ClientHandlers{
		log: log,
	}
}

// Bind binds the client to the server handlers.
func (ch ClientHandlers) Bind(clt *tcp.Client) {
	ch.log.Info(clt.Context(), "client-bind", "userID", clt.UserID())

	clt.Reader = bufio.NewReader(clt.Conn)
}

// Read reads data from the client connection.
func (ch ClientHandlers) Read(clt *tcp.Client) ([]byte, int, error) {
	bufReader := clt.Reader.(*bufio.Reader)

	clt.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	line, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, 0, err
	}

	return []byte(line), len(line), nil
}

// Process processes the request from the client.
func (ch ClientHandlers) Process(r *tcp.Request, clt *tcp.Client) {
}

// Drop is called when a connection is dropped.
func (ch ClientHandlers) Drop(clt *tcp.Client) {
	ch.log.Info(clt.Context(), "client-drop", "userID", clt.UserID())
}

// =============================================================================

// ServerHandlers implements the Handlers interface for the TCP server.
type ServerHandlers struct {
	log      *logger.Logger
	uiCltMgr UIClientManager
}

// NewServerHandlers creates a new instance of ServerHandlers.
func NewServerHandlers(log *logger.Logger, uiCltMgr UIClientManager) *ServerHandlers {
	return &ServerHandlers{
		log:      log,
		uiCltMgr: uiCltMgr,
	}
}

// Bind binds the client to the server handlers.
func (sh ServerHandlers) Bind(clt *tcp.Client) {
	sh.log.Info(clt.Context(), "server-bind", "userID", clt.UserID())

	bufReader := bufio.NewReader(clt.Conn)

	clt.Reader = bufReader

	// -------------------------------------------------------------------------
	// PERFORM HANDSHAKE TO RECEIVE USER ID

	clt.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	line, err := bufReader.ReadString('\n')
	if err != nil {
		sh.log.Info(clt.Context(), "server-bind: handshake-wait", "ERROR", err)
		clt.Conn.Close()
		return
	}

	sh.log.Info(clt.Context(), "server-bind: handshake", "handshake", line)

	var handshake struct {
		UserID string `json:"user_id"`
	}

	if err := json.Unmarshal([]byte(line), &handshake); err != nil {
		sh.log.Info(clt.Context(), "server-bind: handshake-unmarshal", "ERROR", err)
		clt.Conn.Close()
		return
	}

	clt.SetUserID(handshake.UserID)

	// -------------------------------------------------------------------------

	msg := [][]byte{[]byte("EVENT"), []byte("TCP-CONN")}

	from := UIUser{
		ID: common.HexToAddress(clt.UserID()),
	}

	for _, conn := range sh.uiCltMgr.Connections() {
		to := UIUser{
			UIConn: conn.Conn,
		}

		if err := uiSendMessage(from, to, 0, false, msg); err != nil {
			sh.log.Info(clt.Context(), "uilisten: send", "ERROR", err)
		}
	}
}

// Read reads data from the client connection.
func (sh ServerHandlers) Read(clt *tcp.Client) ([]byte, int, error) {
	bufReader := clt.Reader.(*bufio.Reader)

	clt.Conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	line, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, 0, err
	}

	return []byte(line), len(line), nil
}

// Process processes the request from the client.
func (sh ServerHandlers) Process(r *tcp.Request, clt *tcp.Client) {
	ctx := r.Context

	var natsMsg natsInOutMessage
	if err := json.Unmarshal(r.Data, &natsMsg); err != nil {
		sh.log.Info(ctx, "server-process: unmarshal", "ERROR", err)
		return
	}

	sh.log.Info(ctx, "server-process: msg recv", "fromNonce", natsMsg.FromNonce, "from", natsMsg.FromID, "to", natsMsg.ToID, "encrypted", natsMsg.Encrypted, "message", natsMsg.Msg, "fromName", natsMsg.FromName)

	dataThatWasSign := struct {
		ToID      common.Address
		Msg       [][]byte
		FromNonce uint64
	}{
		ToID:      natsMsg.ToID,
		Msg:       natsMsg.Msg,
		FromNonce: natsMsg.FromNonce,
	}

	id, err := signature.FromAddress(dataThatWasSign, natsMsg.V, natsMsg.R, natsMsg.S)
	if err != nil {
		sh.log.Info(ctx, "server-process: fromAddress", "ERROR", err)
		return
	}

	if id != natsMsg.FromID.Hex() {
		sh.log.Info(r.Context, "server-process: signature check", "status", "signature does not match")
		return
	}

	// If the user is found, send the message directly to the user.
	to, err := sh.uiCltMgr.Retrieve(ctx, natsMsg.ToID)
	if err == nil {
		sh.log.Info(ctx, "server-process: msg sent over web socket", "from", natsMsg.FromID, "to", natsMsg.ToID)

		from := UIUser{
			ID:   natsMsg.FromID,
			Name: natsMsg.FromName,
		}

		if err := uiSendMessage(from, to, natsMsg.FromNonce, natsMsg.Encrypted, natsMsg.Msg); err != nil {
			sh.log.Info(ctx, "server-process: send", "ERROR", err)
		}

		return
	}

	if !errors.Is(err, ErrNotExists) {
		sh.log.Info(ctx, "server-process: retrieve", "ERROR", err)
	}

	// We don't have a web socket connection for the user then drop the
	// message on the floor because we can't deliver it.

	sh.log.Info(ctx, "server-process: retrieve", "status", "user not found")
}

// Drop is called when a connection is dropped.
func (sh ServerHandlers) Drop(clt *tcp.Client) {
	sh.log.Info(clt.Context(), "server-drop", "userID", clt.UserID())

	msg := [][]byte{[]byte("EVENT"), []byte("TCP-DROP")}

	from := UIUser{
		ID: common.HexToAddress(clt.UserID()),
	}

	for _, conn := range sh.uiCltMgr.Connections() {
		to := UIUser{
			UIConn: conn.Conn,
		}

		if err := uiSendMessage(from, to, 0, false, msg); err != nil {
			sh.log.Info(clt.Context(), "uilisten: send", "ERROR", err)
		}
	}
}

// =============================================================================

func (b *Business) tcpSendMessage(clt *tcp.Client, from UIUser, inMsg uiIncomingMessage) error {
	natsMsg := natsInOutMessage{
		CapID:             b.capID,
		FromID:            from.ID,
		FromName:          from.Name,
		uiIncomingMessage: inMsg,
	}

	d, err := json.Marshal(natsMsg)
	if err != nil {
		return fmt.Errorf("send nats marshal message: %w", err)
	}

	if _, err := clt.Writer.Write(d); err != nil {
		return fmt.Errorf("send tcp publish: %w", err)
	}

	if _, err := clt.Writer.Write([]byte{'\n'}); err != nil {
		return fmt.Errorf("send tcp publish EOL: %w", err)
	}

	return nil
}
