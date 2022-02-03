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

func main() {
	fmt.Println("Server running in port 8080")
	db := startDatabase("data/data.json")
	if err := http.ListenAndServe("localhost:8080", handler(db)); err != nil {
		log.Fatalf("Error ListenAndServe(): %s", err.Error())
	}
}

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
	result.CurrentUser = models.User{}
	result.Mu = sync.Mutex{}
	return result
}

func updateDatabase(db *models.Database) {
	db.Mu.Lock()
	byteData, err := json.MarshalIndent(db, "", "    ")
	if err != nil {
		fmt.Printf("Failed to marshal data: %s", err.Error())
	}
	if err := ioutil.WriteFile("data/data.json", byteData, 0644); err != nil {
		fmt.Printf("Failed to write data: %s", err.Error())
	}
	db.Mu.Unlock()
}

func handler(db *models.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var transID int
		if r.URL.Path == "/signin" {
			db.Signin(w, r)
		} else if r.URL.Path == "/signup" {
			db.Signup(w, r)
		} else if r.URL.Path == "/transactions" {
			if db.CurrentUser.UserID == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				utils.SendMessageWithBody(w, false, "Please sign in first.")
			} else {
				db.CurrentUser.ProcessTransaction(db, w, r)
			}
		} else if n, _ := fmt.Sscanf(r.URL.Path, "/transactions/%d", &transID); n == 1 {
			if db.CurrentUser.UserID == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				utils.SendMessageWithBody(w, false, "Please sign in first.")
			} else {
				db.CurrentUser.ProcessTransactionID(transID, db, w, r)
			}
		} else {
			w.WriteHeader(http.StatusNotImplemented)
			utils.SendMessage(w, "Invalid URL or request")
		}
		updateDatabase(db)
	}
}
