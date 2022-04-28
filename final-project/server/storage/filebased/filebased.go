package filebased

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/carlodmendoza/go-training/final-project/server/storage"
	"github.com/rs/zerolog/log"
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
		log.Error().Err(err).Msg("Open file error")
		log.Debug().Msg("Using initial sample data")
		file, _ = os.OpenFile("storage/filebased/.data.example", os.O_RDWR|os.O_APPEND, 0644)
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
			log.Error().Err(err).Msg("Read file error")
			return &FilebasedDB{File: file}
		}
		lastScan = tempScan
	}
}

func appendData(fdb *FilebasedDB) error {
	fdb.UserMux.Lock()
	fdb.SessionMux.Lock()
	fdb.CategoryMux.Lock()
	fdb.TransactionMux.Lock()

	err := json.NewEncoder(fdb.File).Encode(fdb)
	if err != nil {
		log.Error().Err(err).Msg("Append data error")
		return err
	}

	fdb.UserMux.Unlock()
	fdb.SessionMux.Unlock()
	fdb.CategoryMux.Unlock()
	fdb.TransactionMux.Unlock()

	return nil
}

func (fdb *FilebasedDB) Shutdown() error {
	return fdb.File.Close()
}
