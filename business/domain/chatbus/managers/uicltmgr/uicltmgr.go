// Package uicltmgr provides a memory based user storage for the chatbus service.
package uicltmgr

import (
	"context"
	"sync"
	"time"

	"github.com/ardanlabs/usdl/business/domain/chatbus"
	"github.com/ardanlabs/usdl/foundation/logger"
	"github.com/ethereum/go-ethereum/common"
)

// UICltMgr provides user management for UI connections.
type UICltMgr struct {
	log     *logger.Logger
	users   map[common.Address]chatbus.UIUser
	muUsers sync.RWMutex
}

// New creates a new manager for UI connections.
func New(log *logger.Logger) *UICltMgr {
	u := UICltMgr{
		log:   log,
		users: make(map[common.Address]chatbus.UIUser),
	}

	return &u
}

// Add adds a new user to the storage.
func (u *UICltMgr) Add(ctx context.Context, usr chatbus.UIUser) error {
	u.muUsers.Lock()
	defer u.muUsers.Unlock()

	if _, exists := u.users[usr.ID]; exists {
		return chatbus.ErrExists
	}

	u.users[usr.ID] = usr

	u.log.Debug(ctx, "chat-adduser", "name", usr.Name, "id", usr.ID)

	return nil
}

// UpdateLastPing updates a user value's ping date/time.
func (u *UICltMgr) UpdateLastPing(ctx context.Context, userID common.Address) error {
	u.muUsers.Lock()
	defer u.muUsers.Unlock()

	usr, exists := u.users[userID]
	if !exists {
		return chatbus.ErrNotExists
	}

	usr.LastPing = time.Now()
	u.users[usr.ID] = usr

	u.log.Debug(ctx, "chat-updping", "name", usr.Name, "id", usr.ID, "lastPing", usr.LastPing)

	return nil
}

// UpdateLastPong updates a user value's pong date/time.
func (u *UICltMgr) UpdateLastPong(ctx context.Context, userID common.Address) (chatbus.UIUser, error) {
	u.muUsers.Lock()
	defer u.muUsers.Unlock()

	usr, exists := u.users[userID]
	if !exists {
		return chatbus.UIUser{}, chatbus.ErrNotExists
	}

	usr.LastPong = time.Now()
	u.users[usr.ID] = usr

	u.log.Debug(ctx, "chat-updpong", "name", usr.Name, "id", usr.ID, "lastPong", usr.LastPong)

	return usr, nil
}

// Remove removes a user from the storage.
func (u *UICltMgr) Remove(ctx context.Context, userID common.Address) {
	u.muUsers.Lock()
	defer u.muUsers.Unlock()

	usr, exists := u.users[userID]
	if !exists {
		u.log.Debug(ctx, "chat-removeuser", "userID", userID, "status", "does not exists")
		return
	}

	delete(u.users, userID)

	u.log.Debug(ctx, "chat-removeuser", "name", usr.Name, "id", usr.ID)
}

// Connections returns all the know users with their connections. A connection
// that is not valid shouldn't be used.
func (u *UICltMgr) Connections() map[common.Address]chatbus.UIConnection {
	u.muUsers.RLock()
	defer u.muUsers.RUnlock()

	m := make(map[common.Address]chatbus.UIConnection)
	for id, usr := range u.users {
		m[id] = chatbus.UIConnection{
			Conn:     usr.UIConn,
			LastPing: usr.LastPing,
			LastPong: usr.LastPong,
		}
	}

	return m
}

// Retrieve retrieves a user from the storage.
func (u *UICltMgr) Retrieve(ctx context.Context, userID common.Address) (chatbus.UIUser, error) {
	u.muUsers.RLock()
	defer u.muUsers.RUnlock()

	usr, exists := u.users[userID]
	if !exists {
		return chatbus.UIUser{}, chatbus.ErrNotExists
	}

	return usr, nil
}
