package models

import (
	"encoding/json"
	"final-project/server/utils"
	"fmt"
	"net/http"
	"sync"
)

type Database struct {
	mu          sync.Mutex
	users       []User
	nextUserID  int
	currentUser User
}

func (db *Database) InitializeDatabase() {
	db.users = []User{}
}

func (db *Database) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/signin" {
			db.signin(w, r)
		} else if r.URL.Path == "/signup" {
			db.signup(w, r)
		} else if r.URL.Path == "/transactions" {
			if db.currentUser.UserID == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				utils.SendMessageWithBody(w, false, "Please login first.")
			} else {
				// VIEW transactions or transaction
				// POST transaction
				// PUT transaction
				// DELETE transaction
				fmt.Println("Handle transactions")
			}
		} else {
			w.WriteHeader(http.StatusNotImplemented)
			utils.SendMessage(w, "Invalid URL or request")
		}
	}
}

func (db *Database) signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		if creds.Username == "" || creds.Password == "" {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
		} else {
			if user := db.authenticateUser(creds); user != nil {
				utils.SendMessageWithBody(w, true, "Logged in successfully!")
				db.currentUser = *user
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				utils.SendMessageWithBody(w, false, "Invalid username or password.")
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

func (db *Database) signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		if user.Name == "" || user.Credentials.Username == "" || user.Credentials.Password == "" {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
		} else {
			if tempUser := db.findUser(user.Credentials); tempUser != nil {
				w.WriteHeader(http.StatusConflict)
				utils.SendMessageWithBody(w, false, "Account already exists.")
			} else {
				db.mu.Lock()
				db.nextUserID++
				user.UserID = db.nextUserID
				db.users = append(db.users, user)
				db.mu.Unlock()

				w.WriteHeader(http.StatusCreated)
				utils.SendMessageWithBody(w, true, "Signed up successfully!")
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

func (db *Database) authenticateUser(creds Credentials) *User {
	db.mu.Lock()
	defer db.mu.Unlock()
	for _, user := range db.users {
		if creds.Username == user.Credentials.Username && creds.Password == user.Credentials.Password {
			return &user
		}
	}
	return nil
}

func (db *Database) findUser(creds Credentials) *User {
	db.mu.Lock()
	defer db.mu.Unlock()
	for _, user := range db.users {
		if creds.Username == user.Credentials.Username {
			return &user
		}
	}
	return nil
}
