package sql_test

import (
	"testing"

	"github.com/ardanlabs/usdl/foundation/client"
	"github.com/ardanlabs/usdl/foundation/client/storage/sql"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	db, err := sql.NewDB(".", common.HexToAddress("0xF"))
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestMyAccount(t *testing.T) {
	db, err := sql.NewDB(".", common.HexToAddress("0xF"))
	assert.NoError(t, err)

	account := db.MyAccount()
	assert.Equal(t, common.HexToAddress("0xF"), account.ID)
	assert.NotEmpty(t, account.Name)

	// If the database is already created, but NewDB is called with a different address,
	// it will update the account ID to the new address. This is because the source
	// of truth for the account is the zarf/client/id folder.
	db, err = sql.NewDB(".", common.HexToAddress("0xFFFFFFFF"))
	assert.NoError(t, err)

	account = db.MyAccount()
	assert.Equal(t, common.HexToAddress("0xFFFFFFFF"), account.ID)
	assert.NotEmpty(t, account.Name)
}

func TestInsertContact(t *testing.T) {
	db, err := sql.NewDB(".", common.HexToAddress("0xF"))
	assert.NoError(t, err)

	err = db.CleanTables()
	assert.NoError(t, err)

	user, err := db.InsertContact(common.HexToAddress("0x1"), "test_user_name")
	assert.NoError(t, err)
	assert.Equal(t, common.HexToAddress("0x1"), user.ID)
	assert.Equal(t, "test_user_name", user.Name)
}

func TestQueryContactByID(t *testing.T) {
	db, err := sql.NewDB(".", common.HexToAddress("0xF"))
	assert.NoError(t, err)

	err = db.CleanTables()
	assert.NoError(t, err)

	user, err := db.InsertContact(common.HexToAddress("0x1"), "test_user_name")
	assert.NoError(t, err)

	user, err = db.QueryContactByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, common.HexToAddress("0x1"), user.ID)
	assert.Equal(t, "test_user_name", user.Name)
}

func TestContacts(t *testing.T) {
	db, err := sql.NewDB(".", common.HexToAddress("0xF"))
	assert.NoError(t, err)

	err = db.CleanTables()
	assert.NoError(t, err)

	contacts := db.Contacts()
	assert.Len(t, contacts, 0)

	user, err := db.InsertContact(common.HexToAddress("0x1"), "test_user_name")
	assert.NoError(t, err)

	err = db.InsertMessage(user.ID, client.Message{
		Name:    "test_message_name",
		Content: [][]byte{[]byte("test_message_content")},
	})
	assert.NoError(t, err)

	user2, err := db.InsertContact(common.HexToAddress("0x2"), "test_user_name_2")
	assert.NoError(t, err)

	contacts = db.Contacts()
	assert.Len(t, contacts, 2)
	assert.Equal(t, user.ID, contacts[0].ID)
	assert.Equal(t, user.Name, contacts[0].Name)
	assert.Equal(t, user.AppLastNonce, contacts[0].AppLastNonce)
	assert.Equal(t, user.LastNonce, contacts[0].LastNonce)
	assert.Equal(t, user.Key, contacts[0].Key)
	assert.Len(t, contacts[0].Messages, 0) // Messages are not loaded in the Contacts() method

	assert.Equal(t, user2.ID, contacts[1].ID)
	assert.Equal(t, user2.Name, contacts[1].Name)
	assert.Equal(t, user2.AppLastNonce, contacts[1].AppLastNonce)
	assert.Equal(t, user2.LastNonce, contacts[1].LastNonce)
	assert.Equal(t, user2.Key, contacts[1].Key)
	assert.Len(t, contacts[1].Messages, 0)
}

func TestInsertMessage(t *testing.T) {
	db, err := sql.NewDB(".", common.HexToAddress("0xF"))
	assert.NoError(t, err)

	err = db.CleanTables()
	assert.NoError(t, err)

	user, err := db.InsertContact(common.HexToAddress("0x1"), "test_user_name")
	assert.NoError(t, err)

	err = db.InsertMessage(user.ID, client.Message{
		Name:    "test_message_name",
		Content: [][]byte{[]byte("test_message_content")},
	})
	assert.NoError(t, err)

	user, err = db.QueryContactByID(user.ID)
	assert.NoError(t, err)
	assert.Len(t, user.Messages, 1)
	assert.Equal(t, "test_message_name", user.Messages[0].Name)
	assert.Equal(t, [][]byte{[]byte("test_message_content")}, user.Messages[0].Content)
	assert.NotEmpty(t, user.Messages[0].DateCreated)
	assert.NotZero(t, user.Messages[0].DateCreated)

	err = db.InsertMessage(user.ID, client.Message{
		Name:    "test_message_name_2",
		Content: [][]byte{[]byte("test_message_content_2_part_1"), []byte("test_message_content_2_part_2")},
	})
	assert.NoError(t, err)

	user, err = db.QueryContactByID(user.ID)
	assert.NoError(t, err)
	assert.Len(t, user.Messages, 2)
	assert.Equal(t, "test_message_name_2", user.Messages[1].Name)
	assert.Equal(t, [][]byte{[]byte("test_message_content_2_part_1"), []byte("test_message_content_2_part_2")}, user.Messages[1].Content)

}

func TestUpdateAppNonce(t *testing.T) {
	db, err := sql.NewDB(".", common.HexToAddress("0xF"))
	assert.NoError(t, err)

	err = db.CleanTables()
	assert.NoError(t, err)

	user, err := db.InsertContact(common.HexToAddress("0x1"), "test_user_name")
	assert.NoError(t, err)

	err = db.UpdateAppNonce(user.ID, 1)
	assert.NoError(t, err)

	user, err = db.QueryContactByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), user.AppLastNonce)
}

func TestUpdateContactNonce(t *testing.T) {
	db, err := sql.NewDB(".", common.HexToAddress("0xF"))
	assert.NoError(t, err)

	err = db.CleanTables()
	assert.NoError(t, err)

	user, err := db.InsertContact(common.HexToAddress("0x1"), "test_user_name")
	assert.NoError(t, err)

	err = db.UpdateContactNonce(user.ID, 1)
	assert.NoError(t, err)

	user, err = db.QueryContactByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), user.LastNonce)

}

func TestUpdateContactKey(t *testing.T) {
	db, err := sql.NewDB(".", common.HexToAddress("0xF"))
	assert.NoError(t, err)

	err = db.CleanTables()
	assert.NoError(t, err)

	user, err := db.InsertContact(common.HexToAddress("0x1"), "test_user_name")
	assert.NoError(t, err)

	err = db.UpdateContactKey(user.ID, "test_key")
	assert.NoError(t, err)

	user, err = db.QueryContactByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "test_key", user.Key)
}

func TestCleanTables(t *testing.T) {
	db, err := sql.NewDB(".", common.HexToAddress("0xF"))
	assert.NoError(t, err)

	err = db.CleanTables()
	assert.NoError(t, err)

	contacts := db.Contacts()
	assert.Len(t, contacts, 0)

	_, err = db.InsertContact(common.HexToAddress("0x1"), "test_user_name")
	assert.NoError(t, err)

	contacts = db.Contacts()
	assert.Len(t, contacts, 1)

	err = db.CleanTables()
	assert.NoError(t, err)

	contacts = db.Contacts()
	assert.Len(t, contacts, 0)
}
