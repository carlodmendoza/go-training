package filebased

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"server/storage"
	"sync"
)

type FilebasedDB struct {
	Users             map[string]*storage.User   `json:"users"`
	Sessions          map[string]storage.Session `json:"sessions"`
	Categories        map[int]storage.Category   `json:"categories"`
	Transactions      []storage.Transaction      `json:"transactions"`
	NextUserID        int                        `json:"nextUserID"`
	NextTransactionID int                        `json:"nextTransactionID"`
	Mu                sync.Mutex
}

var filePath = "../deploy/dev/server/storage/data.json"
var FileDB = startDatabase(filePath)

// startDatabase reads the contents of a json file
// that acts as the database. The result is returned
// as a Database.
func startDatabase(filepath string) *FilebasedDB {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Failed to open json file: %s", err)
	}
	defer file.Close()

	byteData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read json file: %s", err)
	}

	var result *FilebasedDB
	parseErr := json.Unmarshal([]byte(byteData), &result)
	if parseErr != nil {
		log.Fatalf("Failed to parse json file: %s", err)
	}

	return result
}

// updateDatabase writes to a json file that acts as the
// database given a Database.
func updateDatabase(fdb *FilebasedDB) {
	fdb.Mu.Lock()

	byteData, err := json.MarshalIndent(fdb, "", "    ")
	if err != nil {
		fmt.Printf("Failed to marshal data: %s\n", err)
	}

	writeErr := ioutil.WriteFile(filePath, byteData, 0644)
	if writeErr != nil {
		fmt.Printf("Failed to write data: %s\n", err)
	}

	fdb.Mu.Unlock()
}
