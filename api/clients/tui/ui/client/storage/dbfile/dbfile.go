// Package dbfile provides a file based database for the client application.
package dbfile

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/ardanlabs/usdl/api/clients/tui/ui/client"
	"github.com/ethereum/go-ethereum/common"
)

type DB struct {
	myAccount  client.MyAccount
	privKeyRSA *rsa.PrivateKey
	contacts   map[common.Address]client.User
	mu         sync.RWMutex
}

func NewDB(filePath string, id client.ID, jwt string) (*DB, error) {
	df, err := newDB(filePath, id.MyAccountID, jwt)
	if err != nil {
		return nil, fmt.Errorf("newDB: %w", err)
	}

	contacts := make(map[common.Address]client.User, len(df.Contacts))
	for _, usr := range df.Contacts {
		contacts[usr.ID] = client.User{
			ID:           usr.ID,
			Name:         usr.Name,
			AppLastNonce: usr.AppLastNonce,
			LastNonce:    usr.LastNonce,
			Key:          usr.Key,
			TCPHost:      usr.TCPHost,
		}
	}

	db := DB{
		myAccount: client.MyAccount{
			ID:          df.MyAccount.ID,
			Name:        df.MyAccount.Name,
			ProfilePath: df.MyAccount.ProfilePath,
			JWT:         df.MyAccount.JWT,
		},
		privKeyRSA: id.PrivKeyRSA,
		contacts:   contacts,
	}

	return &db, nil
}

func (db *DB) MyAccount() client.MyAccount {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.myAccount
}

func (db *DB) Contacts() []client.User {
	db.mu.RLock()
	defer db.mu.RUnlock()

	users := make([]client.User, 0, len(db.contacts))
	for _, user := range db.contacts {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Name <= users[j].Name
	})

	return users
}

func (db *DB) QueryContactByID(id common.Address) (client.User, error) {
	u, err := func() (client.User, error) {
		db.mu.RLock()
		defer db.mu.RUnlock()

		u, exists := db.contacts[id]
		if !exists {
			return client.User{}, fmt.Errorf("contact not found")
		}

		return u, nil
	}()

	if err != nil {
		return client.User{}, err
	}

	if len(u.Messages) > 0 {
		return u, nil
	}

	// -------------------------------------------------------------------------

	msgs, err := readMsgsFromDisk(id)
	if err != nil {
		return client.User{}, fmt.Errorf("read messages: %w", err)
	}

	messages := make([]client.Message, len(msgs))
	for i, msg := range msgs {
		decryptedData := make([][]byte, len(msg.Content))

		for i, msg := range msg.Content {
			dd, err := rsa.DecryptPKCS1v15(rand.Reader, db.privKeyRSA, []byte(msg))
			if err != nil {
				return client.User{}, fmt.Errorf("decrypting message: %w", err)
			}

			decryptedData[i] = dd
		}

		messages[i] = client.Message{
			Name:        msg.Name,
			Content:     decryptedData,
			DateCreated: msg.DateCreated.Local(),
		}
	}

	// -------------------------------------------------------------------------

	func() {
		db.mu.Lock()
		defer db.mu.Unlock()

		u.Messages = messages
		db.contacts[id] = u
	}()

	return u, nil
}

func (db *DB) InsertContact(id common.Address, name string) (client.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// -------------------------------------------------------------------------
	// Update in the in-memory cache of contacts.

	db.contacts[id] = client.User{
		ID:   id,
		Name: name,
	}

	// -------------------------------------------------------------------------
	// Update the local file.

	df, err := readDBFromDisk()
	if err != nil {
		return client.User{}, fmt.Errorf("config read: %w", err)
	}

	dfu := dataFileUser{
		ID:   id,
		Name: name,
	}

	df.Contacts = append(df.Contacts, dfu)

	flushDBToDisk(df)

	// -------------------------------------------------------------------------
	// Return the new contact.

	u := client.User{
		ID:   id,
		Name: name,
	}

	return u, nil
}

func (db *DB) InsertMessage(id common.Address, msg client.Message) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	u, exists := db.contacts[id]
	if !exists {
		return fmt.Errorf("contact not found")
	}

	u.Messages = append(u.Messages, msg)
	db.contacts[id] = u

	encryptedData := make([][]byte, len(msg.Content))

	for i, msg := range msg.Content {
		ed, err := rsa.EncryptPKCS1v15(rand.Reader, &db.privKeyRSA.PublicKey, msg)
		if err != nil {
			return fmt.Errorf("encrypting message: %w", err)
		}

		encryptedData[i] = ed
	}

	m := message{
		ID:          msg.From,
		Name:        msg.Name,
		Content:     encryptedData,
		DateCreated: time.Now().UTC(),
		Encrypted:   msg.Encrypted,
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
