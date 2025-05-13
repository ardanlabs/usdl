package dbfile

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const (
	dbDirName     = "db"
	dbMsgsDirName = "msgs"
	dbFileName    = "data.json"
)

var (
	dbFileDir string
	dbMsgsDir string
	dbFile    string
)

type message struct {
	ID          common.Address `json:"id"`
	Name        string         `json:"name"`
	Encrypted   bool           `json:"encrypted"`
	Content     [][]byte       `json:"content"`
	DateCreated time.Time      `json:"date_created"`
}

type myAccount struct {
	ID   common.Address `json:"id"`
	Name string         `json:"name"`
}

type dataFileUser struct {
	ID           common.Address `json:"id"`
	Name         string         `json:"name"`
	AppLastNonce uint64         `json:"app_last_nonce"`
	LastNonce    uint64         `json:"last_nonce"`
	Key          string         `json:"key,omitempty"`
}

type dataFile struct {
	MyAccount myAccount      `json:"my_account"`
	Contacts  []dataFileUser `json:"contacts"`
}

func newDB(filePath string, myAccountID common.Address) (dataFile, error) {
	dbFileDir = filepath.Join(filePath, dbDirName)
	dbMsgsDir = filepath.Join(filePath, dbDirName, dbMsgsDirName)
	dbFile = filepath.Join(filePath, dbDirName, dbFileName)

	os.MkdirAll(dbFileDir, os.ModePerm)
	os.MkdirAll(dbMsgsDir, os.ModePerm)

	var df dataFile

	_, err := os.Stat(dbFile)
	switch {
	case err != nil:
		df, err = createDBOnDisk(myAccountID)

	default:
		df, err = readDBFromDisk()
		if df.MyAccount.ID != myAccountID {
			return dataFile{}, fmt.Errorf("id mismatch: got: %s exp: %s", df.MyAccount.ID.Hex(), myAccountID.Hex())
		}
	}

	if err != nil {
		return dataFile{}, fmt.Errorf("config: %w", err)
	}

	return df, nil
}

func createDBOnDisk(myAccountID common.Address) (dataFile, error) {
	f, err := os.Create(dbFile)
	if err != nil {
		return dataFile{}, fmt.Errorf("config data file create: %w", err)
	}
	defer f.Close()

	df := dataFile{
		MyAccount: myAccount{
			ID:   myAccountID,
			Name: "Anonymous",
		},
		Contacts: []dataFileUser{
			{
				ID:   common.Address{},
				Name: "Sample Contact",
			},
		},
	}

	jsonDoc, err := json.MarshalIndent(df, "", "    ")
	if err != nil {
		return dataFile{}, fmt.Errorf("config data file marshal: %w", err)
	}

	if _, err := f.Write(jsonDoc); err != nil {
		return dataFile{}, fmt.Errorf("config data file write: %w", err)
	}

	return df, nil
}

func readDBFromDisk() (dataFile, error) {
	f, err := os.Open(dbFile)
	if err != nil {
		return dataFile{}, fmt.Errorf("id data file open: %w", err)
	}
	defer f.Close()

	var df dataFile
	if err := json.NewDecoder(f).Decode(&df); err != nil {
		return dataFile{}, fmt.Errorf("config decode: %w", err)
	}

	return df, nil
}

func flushDBToDisk(df dataFile) error {
	f, err := os.Create(dbFile)
	if err != nil {
		return fmt.Errorf("config data file create: %w", err)
	}
	defer f.Close()

	jsonDF, err := json.MarshalIndent(df, "", "    ")
	if err != nil {
		return fmt.Errorf("config data file marshal: %w", err)
	}

	if _, err := f.Write(jsonDF); err != nil {
		return fmt.Errorf("config data file write: %w", err)
	}

	return nil
}

func readMsgsFromDisk(id common.Address) ([]message, error) {
	fileName := filepath.Join(dbMsgsDir, id.Hex()+".msg")

	f, err := os.Open(fileName)
	if err != nil {
		return nil, nil
	}
	defer f.Close()

	var msgs []message

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var msg message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			return nil, fmt.Errorf("unmarshal: %w", err)
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func flushMsgToDisk(id common.Address, msg message) error {
	fileName := filepath.Join(dbMsgsDir, id.Hex()+".msg")

	var f *os.File

	_, err := os.Stat(fileName)
	switch {
	case err != nil:

		f, err = os.Create(fileName)
		if err != nil {
			return fmt.Errorf("message file create: %w", err)
		}

	default:
		f, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("message file open: %w", err)
		}
	}

	defer f.Close()

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("message marshal: %w", err)
	}

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("message file write: %w", err)
	}

	if _, err := f.WriteString("\n"); err != nil {
		return fmt.Errorf("message file write: %w", err)
	}

	return nil
}
