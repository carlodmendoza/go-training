package filebased

import (
	"encoding/json"
	"log"
	"os"
	"server/storage"
	"sync"
)

type FilebasedDB struct {
	Users             map[string]storage.User     `json:"users"`
	Sessions          map[string]storage.Session  `json:"sessions"`
	Categories        map[int]storage.Category    `json:"categories"`
	Transactions      map[int]storage.Transaction `json:"transactions"`
	NextUserID        int                         `json:"nextUserID"`
	NextTransactionID int                         `json:"nextTransactionID"`
	UserMux           sync.RWMutex
	SessionMux        sync.RWMutex
	CategoryMux       sync.RWMutex
	TransactionMux    sync.RWMutex
}

var (
	filePath = "../deploy/dev/server/storage/data.json"
	filePtr  = openFile(filePath)
	FileDB   = readFile(filePtr)
)

func openFile(filepath string) *os.File {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Failed to open json file: %s", err)
	}
	defer file.Close()

	return file
}

func readFile(file *os.File) *FilebasedDB {
	var result *FilebasedDB
	err := json.NewDecoder(file).Decode(&result)
	if err != nil {
		log.Fatalf("Failed to decode json file: %s", err.Error())
	}

	return result
}

func writeToFile(file *os.File, fdb *FilebasedDB) {
	fdb.UserMux.Lock()
	fdb.SessionMux.Lock()
	fdb.CategoryMux.Lock()
	fdb.TransactionMux.Lock()

	err := json.NewEncoder(file).Encode(fdb)
	if err != nil {
		log.Fatalf("Failed to encode json file: %s", err.Error())
	}

	fdb.UserMux.Unlock()
	fdb.SessionMux.Unlock()
	fdb.CategoryMux.Unlock()
	fdb.TransactionMux.Unlock()
}
