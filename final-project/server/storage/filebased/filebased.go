package filebased

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"server/storage"
	"sync"
)

type FilebasedDB struct {
	File *os.File

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

func Initialize(filepath string) *FilebasedDB {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	decoder := json.NewDecoder(file)
	var lastScan *FilebasedDB

	for {
		var tempScan *FilebasedDB
		err := decoder.Decode(&tempScan)
		switch {
		case err == io.EOF:
			lastScan.File = file
			return lastScan
		case err != nil:
			log.Fatalf("Failed to read file: %s", err)
		}
		lastScan = tempScan
	}
}

func appendData(fdb *FilebasedDB) {
	fdb.UserMux.Lock()
	fdb.SessionMux.Lock()
	fdb.CategoryMux.Lock()
	fdb.TransactionMux.Lock()

	err := json.NewEncoder(fdb.File).Encode(fdb)
	if err != nil {
		log.Fatalf("Failed to append data: %s", err)
	}

	fdb.UserMux.Unlock()
	fdb.SessionMux.Unlock()
	fdb.CategoryMux.Unlock()
	fdb.TransactionMux.Unlock()
}

func (fdb *FilebasedDB) Shutdown() error {
	return fdb.File.Close()
}
