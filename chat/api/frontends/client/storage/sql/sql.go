// Package sql provides a SQLite database for the chat application.
package sql

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ardanlabs/usdl/chat/api/frontends/client/app"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	dbDirName  = "db"
	dbFileName = "data.db"
)

type DB struct {
	db *gorm.DB
}

type myAccount struct {
	Singleton bool      `gorm:"primaryKey;default:true"`
	ID        string    `gorm:"column:id"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type user struct {
	ID           string    `gorm:"primaryKey;column:id"`
	Name         string    `gorm:"column:name"`
	AppLastNonce uint64    `gorm:"column:app_last_nonce"`
	LastNonce    uint64    `gorm:"column:last_nonce"`
	Key          string    `gorm:"column:key"`
	Messages     []message `gorm:"foreignKey:UserID;column:messages"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

type message struct {
	ID        uint64    `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"column:name"`
	Content   []byte    `gorm:"column:content"`
	UserID    string    `gorm:"column:user_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func NewDB(filePath string, myAccountID common.Address) (*DB, error) {
	dbFileDir := filepath.Join(filePath, dbDirName)
	os.MkdirAll(dbFileDir, os.ModePerm)

	fileName := filepath.Join(dbFileDir, dbFileName)
	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}

	if err := db.AutoMigrate(&user{}, &message{}, myAccount{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	if err := saveMyAccount(db, myAccountID); err != nil {
		return nil, fmt.Errorf("save my account: %w", err)
	}

	return &DB{db: db}, nil
}

func (db *DB) MyAccount() app.MyAccount {
	myAccount := myAccount{
		Singleton: true,
	}
	res := db.db.First(&myAccount)
	if res.Error != nil {
		return app.MyAccount{} // maybe better to return an error
	}

	return app.MyAccount{
		ID:   common.HexToAddress(myAccount.ID),
		Name: myAccount.Name,
	}
}

func (db *DB) InsertContact(id common.Address, name string) (app.User, error) {
	res := db.db.Create(&user{
		ID:   strings.ToLower(id.Hex()),
		Name: name,
	})
	if res.Error != nil {
		return app.User{}, fmt.Errorf("insert contact: %w", res.Error)
	}

	return app.User{
		ID:   id,
		Name: name,
	}, nil
}

func (db *DB) QueryContactByID(id common.Address) (app.User, error) {
	var user user
	if err := db.db.Preload("Messages").Where("LOWER(id) = LOWER(?)", id.Hex()).First(&user).Error; err != nil {
		return app.User{}, fmt.Errorf("query contact: %w", err)
	}

	msgs := make([]app.Message, len(user.Messages))
	for i, msg := range user.Messages {
		msgs[i] = app.Message{
			Name:        msg.Name,
			Content:     msg.Content,
			DateCreated: msg.CreatedAt.Local(),
		}
	}

	return app.User{
		ID:           common.HexToAddress(user.ID),
		Name:         user.Name,
		AppLastNonce: user.AppLastNonce,
		LastNonce:    user.LastNonce,
		Key:          user.Key,
		Messages:     msgs,
	}, nil
}

func (db *DB) Contacts() []app.User {
	var users []user
	res := db.db.Preload("Messages").Find(&users)
	if res.Error != nil {
		return nil // maybe better to return an error
	}

	contacts := make([]app.User, len(users))
	for i, user := range users {
		msgs := make([]app.Message, len(user.Messages))
		for j, msg := range user.Messages {
			msgs[j] = app.Message{
				Name:        msg.Name,
				Content:     msg.Content,
				DateCreated: msg.CreatedAt.Local(),
			}
		}
		contacts[i] = app.User{
			ID:           common.HexToAddress(user.ID),
			Name:         user.Name,
			AppLastNonce: user.AppLastNonce,
			LastNonce:    user.LastNonce,
			Key:          user.Key,
			Messages:     msgs,
		}
	}
	return contacts
}

func (db *DB) InsertMessage(id common.Address, msg app.Message) error {
	res := db.db.Create(&message{
		Name:    msg.Name,
		Content: msg.Content,
		UserID:  strings.ToLower(id.Hex()),
	})
	if res.Error != nil {
		return fmt.Errorf("insert message: %w", res.Error)
	}

	return nil
}

func (db *DB) UpdateAppNonce(id common.Address, nonce uint64) error {
	res := db.db.Model(&user{}).Where("LOWER(id) = LOWER(?)", id.Hex()).Update("app_last_nonce", nonce)
	if res.Error != nil {
		return fmt.Errorf("update app nonce: %w", res.Error)
	}

	return nil
}

func (db *DB) UpdateContactNonce(id common.Address, nonce uint64) error {
	res := db.db.Model(&user{}).Where("LOWER(id) = LOWER(?)", id.Hex()).Update("last_nonce", nonce)
	if res.Error != nil {
		return fmt.Errorf("update contact nonce: %w", res.Error)
	}

	return nil
}

func (db *DB) UpdateContactKey(id common.Address, key string) error {
	res := db.db.Model(&user{}).Where("LOWER(id) = LOWER(?)", id.Hex()).Update("key", key)
	if res.Error != nil {
		return fmt.Errorf("update contact key: %w", res.Error)
	}

	return nil
}

func (db *DB) CleanTables() error {
	if err := db.db.Migrator().DropTable(&user{}, &message{}); err != nil {
		return fmt.Errorf("drop table: %w", err)
	}

	if err := db.db.AutoMigrate(&user{}, &message{}); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	return nil
}

func saveMyAccount(db *gorm.DB, myAccountID common.Address) error {
	var myAcc myAccount
	res := db.First(&myAcc)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("query my account: %w", res.Error)
	}

	if res.Error == gorm.ErrRecordNotFound {
		return db.Save(&myAccount{
			Singleton: true,
			ID:        strings.ToLower(myAccountID.Hex()),
			Name:      "Anonymous",
		}).Error
	}

	myAcc.Singleton = true
	myAcc.ID = strings.ToLower(myAccountID.Hex())
	return db.Save(&myAcc).Error
}
