package app

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type UIScreenWrite func(name string, msg string)

// =============================================================================

type inMessage struct {
	ToID string `json:"toID"`
	Msg  string `json:"msg"`
}

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type outMessage struct {
	From user   `json:"from"`
	Msg  string `json:"msg"`
}

// =============================================================================

type Client struct {
	id   string
	url  string
	conn *websocket.Conn
	cfg  *Config
}

func NewClient(id string, url string, cfg *Config) *Client {
	return &Client{
		id:  id,
		url: url,
		cfg: cfg,
	}
}

func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

func (c *Client) Handshake(name string, write UIScreenWrite) error {
	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	c.conn = conn

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
		ID   string
		Name string
	}{
		ID:   c.id,
		Name: name,
	}

	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	// -------------------------------------------------------------------------

	_, msg, err = conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	// -------------------------------------------------------------------------

	go func() {
		for {
			_, msg, err = conn.ReadMessage()
			if err != nil {
				write("system", fmt.Sprintf("read: %s", err))
				return
			}

			var outMsg outMessage
			if err := json.Unmarshal(msg, &outMsg); err != nil {
				write("system", fmt.Sprintf("unmarshal: %s", err))
				return
			}

			user, err := c.cfg.LookupContact(outMsg.From.ID)
			if err == nil {
				outMsg.From.Name = user.Name
			}

			write(outMsg.From.Name, outMsg.Msg)
		}
	}()

	return nil
}

func (c *Client) Send(to string, msg string) error {
	if c.conn == nil {
		return fmt.Errorf("no connection")
	}

	inMsg := inMessage{
		ToID: to,
		Msg:  msg,
	}

	data, err := json.Marshal(inMsg)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
