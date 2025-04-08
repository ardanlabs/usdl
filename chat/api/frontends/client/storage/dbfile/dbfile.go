package dbfile

import (
	"fmt"
	"sync"

	"github.com/ardanlabs/usdl/chat/api/frontends/client/app"
	"github.com/ethereum/go-ethereum/common"
)

type DB struct {
	myAccount app.MyAccount
	contacts  map[common.Address]app.User
	mu        sync.RWMutex
}

func NewDB(filePath string, myAccountID common.Address) (*DB, error) {
	df, err := newDB(filePath, myAccountID)
	if err != nil {
		return nil, fmt.Errorf("newDB: %w", err)
	}

	contacts := make(map[common.Address]app.User, len(df.Contacts))
	for _, usr := range df.Contacts {
		contacts[usr.ID] = app.User{
			ID:           usr.ID,
			Name:         usr.Name,
			AppLastNonce: usr.AppLastNonce,
			LastNonce:    usr.LastNonce,
			Key:          usr.Key,
		}
	}

	db := DB{
		myAccount: app.MyAccount{
			ID:   df.MyAccount.ID,
			Name: df.MyAccount.Name,
		},
		contacts: contacts,
	}

	return &db, nil
}

func (c *DB) MyAccount() app.MyAccount {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.myAccount
}

func (c *DB) Contacts() []app.User {
	c.mu.RLock()
	defer c.mu.RUnlock()

	users := make([]app.User, 0, len(c.contacts))
	for _, user := range c.contacts {
		users = append(users, user)
	}

	return users
}

func (db *DB) QueryContactByID(id common.Address) (app.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	u, exists := db.contacts[id]
	if !exists {
		return app.User{}, fmt.Errorf("contact not found")
	}

	if len(u.Messages) == 0 {
		msgs, err := readMsgsFromDisk(id)
		if err != nil {
			return app.User{}, fmt.Errorf("read messages: %w", err)
		}

		u.Messages = make([]app.Message, len(msgs))
		for i, msg := range msgs {
			u.Messages[i] = app.Message{
				Name:    msg.Name,
				Content: msg.Content,
			}
		}

		db.contacts[id] = u
	}

	return u, nil
}

func (db *DB) InsertContact(id common.Address, name string) (app.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// -------------------------------------------------------------------------
	// Update in the in-memory cache of contacts.

	db.contacts[id] = app.User{
		ID:   id,
		Name: name,
	}

	// -------------------------------------------------------------------------
	// Update the local file.

	df, err := readDBFromDisk()
	if err != nil {
		return app.User{}, fmt.Errorf("config read: %w", err)
	}

	dfu := dataFileUser{
		ID:   id,
		Name: name,
	}

	df.Contacts = append(df.Contacts, dfu)

	flushDBToDisk(df)

	// -------------------------------------------------------------------------
	// Return the new contact.

	u := app.User{
		ID:   id,
		Name: name,
	}

	return u, nil
}

func (db *DB) InsertMessage(id common.Address, msg app.Message) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	u, exists := db.contacts[id]
	if !exists {
		return fmt.Errorf("contact not found")
	}

	u.Messages = append(u.Messages, msg)
	db.contacts[id] = u

	m := message{
		Name:    msg.Name,
		Content: msg.Content,
	}

	if err := flushMsgToDisk(id, m); err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	return nil
}

func (db *DB) UpdateAppNonce(id common.Address, nonce uint64) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// -------------------------------------------------------------------------
	// Update in the in-memory cache of contacts.

	u, exists := db.contacts[id]
	if !exists {
		return fmt.Errorf("contact not found")
	}

	u.AppLastNonce = nonce

	db.contacts[id] = u

	// -------------------------------------------------------------------------
	// Update the local file.

	df, err := readDBFromDisk()
	if err != nil {
		return fmt.Errorf("config read: %w", err)
	}

	for i, contact := range df.Contacts {
		if contact.ID == id {
			df.Contacts[i].AppLastNonce = nonce
			break
		}
	}

	flushDBToDisk(df)

	return nil
}

func (db *DB) UpdateContactNonce(id common.Address, nonce uint64) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// -------------------------------------------------------------------------
	// Update in the in-memory cache of contacts.

	u, exists := db.contacts[id]
	if !exists {
		return fmt.Errorf("contact not found")
	}

	u.LastNonce = nonce

	db.contacts[id] = u

	// -------------------------------------------------------------------------
	// Update the local file.

	df, err := readDBFromDisk()
	if err != nil {
		return fmt.Errorf("config read: %w", err)
	}

	for i, contact := range df.Contacts {
		if contact.ID == id {
			df.Contacts[i].LastNonce = nonce
			break
		}
	}

	flushDBToDisk(df)

	return nil
}

func (db *DB) UpdateContactKey(id common.Address, key string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// -------------------------------------------------------------------------
	// Update in the in-memory cache of contacts.

	u, exists := db.contacts[id]
	if !exists {
		return fmt.Errorf("contact not found")
	}

	u.Key = key

	db.contacts[id] = u

	// -------------------------------------------------------------------------
	// Update the local file.

	df, err := readDBFromDisk()
	if err != nil {
		return fmt.Errorf("config read: %w", err)
	}

	for i, contact := range df.Contacts {
		if contact.ID == id {
			df.Contacts[i].Key = key
			break
		}
	}

	flushDBToDisk(df)

	return nil
}
