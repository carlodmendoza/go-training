package main

import (
	"encoding/json"
	"final-project/server/models"
	"final-project/server/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

/*
	Main program for running the server, handling requests,
	and starting and reading the file-based storage.
	Author: Carlo Mendoza
*/

func main() {
	fmt.Println("Server running on port 8080")
	db := startDatabase("data/data.json")
	if err := http.ListenAndServe(":8080", handler(db)); err != nil {
		log.Fatalf("Error ListenAndServe(): %s", err.Error())
	}
}

// startDatabase reads the contents of a json file
// that acts as the database. The result is returned
// as a Database.
func startDatabase(filepath string) *models.Database {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Failed to open json file: %s", err.Error())
	}
	defer file.Close()
	byteData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read json file: %s", err.Error())
	}
	var result *models.Database
	if err := json.Unmarshal([]byte(byteData), &result); err != nil {
		log.Fatalf("Failed to parse json file: %s", err.Error())
	}
	result.Mu = sync.Mutex{}
	return result
}

// updateDatabase writes to a json file that acts as the
// database given a Database.
func updateDatabase(db *models.Database) {
	db.Mu.Lock()
	byteData, err := json.MarshalIndent(db, "", "    ")
	if err != nil {
		fmt.Printf("Failed to marshal data: %s\n", err.Error())
	}
	if err := ioutil.WriteFile("data/data.json", byteData, 0644); err != nil {
		fmt.Printf("Failed to write data: %s\n", err.Error())
	}
	db.Mu.Unlock()
}

// handler handles requests to the server depending on the
// request URL given a Database. It updates the database
// every client request. It also authorizes a user to make
// requests given a request token.
func handler(db *models.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var transID int
		if r.URL.Path == "/signin" {
			db.Signin(w, r)
		} else if r.URL.Path == "/signup" {
			db.Signup(w, r)
		} else if r.URL.Path == "/transactions" {
			uid := db.FindUidByToken(r)
			if uid == -1 || uid == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				utils.SendMessageWithBody(w, false, "Unauthorized login.")
			} else {
				db.ProcessTransaction(w, r, uid)
			}
		} else if n, _ := fmt.Sscanf(r.URL.Path, "/transactions/%d", &transID); n == 1 {
			uid := db.FindUidByToken(r)
			if uid == -1 || uid == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				utils.SendMessageWithBody(w, false, "Unauthorized login.")
			} else {
				db.ProcessTransactionID(w, r, uid, transID)
			}
		} else if r.URL.Path == "/categories" {
			uid := db.FindUidByToken(r)
			if uid == -1 || uid == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				utils.SendMessageWithBody(w, false, "Unauthorized login.")
			} else {
				db.ProcessCategories(w, r)
			}
		} else {
			w.WriteHeader(http.StatusNotImplemented)
			utils.SendMessage(w, "Invalid URL or request")
		}
		updateDatabase(db)
	}
}
