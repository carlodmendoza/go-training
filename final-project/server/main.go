package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type User struct {
	UserID            int
	Name              string
	Credentials       UserCredentials
	wallets           []Wallet
	transactions      []Transaction
	nextWalletID      int
	nextTransactionID int
}

type UserCredentials struct {
	Username string
	Password string
}

type Wallet struct {
	WalletID       int
	Name           string
	Currency       string
	InitialBalance float64
}

type Transaction struct {
	TransactionID int
	WalletID      int
	CategoryID    int
	Amount        float64
	Date          string
	Notes         string
}

type Category struct {
	CategoryID int
	Name       string
	Type       string
}

type Database struct {
	nextUserID  int
	mu          sync.Mutex
	users       []User
	currentUser User
}

func main() {
	fmt.Println("Server running in port 8080")
	db := &Database{users: []User{}}
	if err := http.ListenAndServe("localhost:8080", db.handler()); err != nil {
		log.Fatalf("Error ListenAndServe(): %s", err.Error())
	}
}

// handler handles requests to the server
// depending on the request URL
func (db *Database) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/signin" {
			db.signin(w, r)
		} else if r.URL.Path == "/signup" {
			db.signup(w, r)
		} else if r.URL.Path == "/wallets" {
			if db.currentUser.UserID == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				sendMessageWithBody(w, false, "Please login first.")
			} else {
				// GET wallets or wallet
				// POST wallet
				// PUT wallet
				// DELETE wallet
			}
		} else if r.URL.Path == "/transactions" {
			if db.currentUser.UserID == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				sendMessageWithBody(w, false, "Please login first.")
			} else {
				// VIEW transactions or transaction
				// POST transaction
				// PUT transaction
				// DELETE transaction
			}
		} else {
			w.WriteHeader(http.StatusNotImplemented)
			sendMessage(w, "Invalid URL or request")
		}
	}
}

func (db *Database) signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var creds UserCredentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			sendMessage(w, "400 Bad Request")
			return
		}
		if creds.Username == "" || creds.Password == "" {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			w.WriteHeader(http.StatusBadRequest)
			sendMessage(w, "400 Bad Request")
		} else {
			if user := db.authenticateUser(creds); user != nil {
				sendMessageWithBody(w, true, "Logged in successfully!")
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				sendMessageWithBody(w, false, "Invalid username or password.")
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		sendMessage(w, "405 Method not allowed")
	}
}

func (db *Database) signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			sendMessage(w, "400 Bad Request")
			return
		}
		if user.Name == "" || user.Credentials.Username == "" || user.Credentials.Password == "" {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			w.WriteHeader(http.StatusBadRequest)
			sendMessage(w, "400 Bad Request")
		} else {
			if tempUser := db.findUser(user.Credentials); tempUser != nil {
				w.WriteHeader(http.StatusConflict)
				sendMessageWithBody(w, false, "Account already exists.")
			} else {
				db.mu.Lock()
				db.nextUserID++
				user.UserID = db.nextUserID
				db.users = append(db.users, user)
				db.mu.Unlock()

				w.WriteHeader(http.StatusCreated)
				sendMessageWithBody(w, true, "Signed up successfully!")
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		sendMessage(w, "405 Method not allowed")
	}
}

func (db *Database) authenticateUser(creds UserCredentials) *User {
	db.mu.Lock()
	defer db.mu.Unlock()
	for _, user := range db.users {
		if creds.Username == user.Credentials.Username && creds.Password == user.Credentials.Password {
			return &user
		}
	}
	return nil
}

func (db *Database) findUser(creds UserCredentials) *User {
	db.mu.Lock()
	defer db.mu.Unlock()
	for _, user := range db.users {
		if creds.Username == user.Credentials.Username {
			return &user
		}
	}
	return nil
}
