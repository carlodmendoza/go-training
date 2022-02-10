package models

import (
	"encoding/json"
	"final-project/server/utils"
	"fmt"
	"net/http"
	"sync"
	"time"
)

/*
	Database program contains all fields and methods
	of Database that make it possible to process
	requests from a client.
	Author: Carlo Mendoza
*/

type Database struct {
	Users             []User        `json:"users"`
	Sessions          []Session     `json:"sessions"`
	Credentials       []Credentials `json:"credentials"`
	Transactions      []Transaction `json:"transactions"`
	Categories        []Category    `json:"categories"`
	NextUserID        int           `json:"nextUserID"`
	NextTransactionID int           `json:"nextTransactionID"`
	Mu                sync.Mutex
}

// Signin handles a sign in request by a client.
// Upon successful sign in, a generated token
// is given as a cookie to client for authorizing
// future requests.
func (db *Database) Signin(w http.ResponseWriter, r *http.Request) {
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
			if userID, ok := db.authenticateUser(creds); ok {
				var session Session
				session.Token = utils.GenerateSessionToken()
				session.Timestamp = time.Now().Format("2006-01-02 15:04:05")
				session.UserID = userID
				// if existing session, replace it with a new token and timestamp
				if _, index, ok := db.findSessionByUid(userID); ok {
					db.Mu.Lock()
					db.Sessions[index] = session
					db.Mu.Unlock()
				} else {
					db.Mu.Lock()
					db.Sessions = append(db.Sessions, session)
					db.Mu.Unlock()
				}
				// send token as Cookie to client
				cookie := http.Cookie{Name: "Token", Value: session.Token}
				http.SetCookie(w, &cookie)
				utils.SendMessageWithBody(w, true, "Logged in successfully!")
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

// Signup handles a sign up request by a client.
// It checks if an account already exists.
func (db *Database) Signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var user User
		var creds Credentials
		var reqBody map[string]string
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		if reqBody["name"] == "" || reqBody["username"] == "" || reqBody["password"] == "" {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
		} else {
			if tempUser := db.findCredentialsByUsername(reqBody["username"]); tempUser != nil {
				w.WriteHeader(http.StatusConflict)
				utils.SendMessageWithBody(w, false, "Account already exists.")
			} else {
				db.Mu.Lock()
				db.NextUserID++
				user.UserID = db.NextUserID
				user.Name = reqBody["name"]
				db.Users = append(db.Users, user)

				creds.Username = reqBody["username"]
				creds.Password = reqBody["password"]
				creds.UserID = db.NextUserID
				db.Credentials = append(db.Credentials, creds)
				db.Mu.Unlock()

				w.WriteHeader(http.StatusCreated)
				utils.SendMessageWithBody(w, true, "Signed up successfully!")
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

// ProcessTransaction handles a transaction/ request by a client
// given a user ID. The client can either get all transactions,
// add new transaction, or delete all transactions.
func (db *Database) ProcessTransaction(w http.ResponseWriter, r *http.Request, userID int) {
	switch r.Method {
	case "GET":
		transactions := db.findTransactionsByUid(userID)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(transactions); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			utils.SendMessageWithBody(w, false, "500 Internal Server Error")
			return
		}
	case "POST":
		var transaction Transaction
		if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		if transaction.Amount == 0 || transaction.Date == "" || transaction.CategoryID == 0 {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
		} else {
			if ok := db.validateNewTransaction(w, r, transaction); ok {
				db.Mu.Lock()
				db.NextTransactionID++
				transaction.TransactionID = db.NextTransactionID
				transaction.UserID = userID
				db.Transactions = append(db.Transactions, transaction)
				db.Mu.Unlock()

				w.WriteHeader(http.StatusCreated)
				utils.SendMessageWithBody(w, true, "Transaction added successfully!")
			}
		}
	case "DELETE":
		transactions := db.findTransactionsByUid(userID)
		if len(transactions) == 0 {
			w.WriteHeader(http.StatusNotFound)
			utils.SendMessageWithBody(w, false, "No transactions found.")
		} else {
			db.deleteAllTransactionsByUid(userID)
			utils.SendMessageWithBody(w, true, "All transactions deleted successfully.")
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

// ProcessTransactionID handles a transaction/id request by a client
// given a user ID and a transaction ID. The client can either get,
// update, or delete a transaction.
func (db *Database) ProcessTransactionID(w http.ResponseWriter, r *http.Request, userID, transID int) {
	switch r.Method {
	case "GET":
		if trans, _, ok := db.findTransactionByTid(userID, transID); ok {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(&trans); err != nil {
				fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				utils.SendMessageWithBody(w, false, "500 Internal Server Error")
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			utils.SendMessageWithBody(w, false, "Transaction not found.")
		}
	case "PUT":
		if trans, index, ok := db.findTransactionByTid(userID, transID); ok {
			transaction := *trans
			if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
				fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
				w.WriteHeader(http.StatusBadRequest)
				utils.SendMessage(w, "400 Bad Request")
				return
			}
			if ok := db.validateNewTransaction(w, r, transaction); ok {
				db.Mu.Lock()
				// ensures that original IDs are not changed
				transaction.TransactionID = trans.TransactionID
				transaction.UserID = trans.UserID
				db.Transactions[index] = transaction
				db.Mu.Unlock()

				utils.SendMessageWithBody(w, true, "Transaction updated successfully!")
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			utils.SendMessageWithBody(w, false, "Transaction not found.")
		}
	case "DELETE":
		if _, index, ok := db.findTransactionByTid(userID, transID); ok {
			db.Mu.Lock()
			db.Transactions = append(db.Transactions[:index], db.Transactions[index+1:]...)
			db.Mu.Unlock()
			utils.SendMessageWithBody(w, true, "Transaction deleted successfully.")
		} else {
			w.WriteHeader(http.StatusNotFound)
			utils.SendMessageWithBody(w, false, "Transaction not found.")
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

// ProcessCategories handles a categories/ request by a client.
// The client can get all categories.
func (db *Database) ProcessCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(db.Categories); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			utils.SendMessageWithBody(w, false, "500 Internal Server Error")
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

// FindUidByToken returns the user ID given a token
// sent from a request cookie.
func (db *Database) FindUidByToken(r *http.Request) int {
	tokenCookie, err := r.Cookie("Token")
	if err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
		return -1
	}
	db.Mu.Lock()
	defer db.Mu.Unlock()
	for _, sesh := range db.Sessions {
		if tokenCookie.Value == sesh.Token {
			return sesh.UserID
		}
	}
	return 0
}

// authenticateUser returns the user ID given a username.
// If username doesn't exist, it returns 0 and false.
func (db *Database) authenticateUser(creds Credentials) (int, bool) {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	for _, userCreds := range db.Credentials {
		if creds.Username == userCreds.Username && creds.Password == userCreds.Password {
			return userCreds.UserID, true
		}
	}
	return 0, false
}

// findSessionByUid returns a Session pointer and index given a user ID.
// It is used for checking existing sessions. If a session doesn't
// exist, it returns nil, false index, and false.
func (db *Database) findSessionByUid(uid int) (*Session, int, bool) {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	for i, sesh := range db.Sessions {
		if uid == sesh.UserID {
			return &sesh, i, true
		}
	}
	return nil, -1, false
}

// findCredentialsByUsername returns a Credentials pointer given
// a username. It is used for checking existing accounts. If
// credentials don't exist, it returns a nil pointer.
func (db *Database) findCredentialsByUsername(username string) *Credentials {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	for _, user := range db.Credentials {
		if username == user.Username {
			return &user
		}
	}
	return nil
}

// findTransactionsByUid returns a list of Transaction given
// a user ID. If there are no existing user transactions,
// it returns an empty list.
func (db *Database) findTransactionsByUid(uid int) []Transaction {
	transactions := []Transaction{}
	db.Mu.Lock()
	defer db.Mu.Unlock()
	for _, trans := range db.Transactions {
		if uid == trans.UserID {
			transactions = append(transactions, trans)
		}
	}
	return transactions
}

// findTransactionByTid returns a Transaction pointer and index given a user ID.
// It is used for checking existing transactions. If a transaction doesn't
// exist, it returns nil, false index, and false.
func (db *Database) findTransactionByTid(uid, tid int) (*Transaction, int, bool) {
	transactions := db.findTransactionsByUid(uid)
	db.Mu.Lock()
	defer db.Mu.Unlock()
	for i, trans := range transactions {
		if tid == trans.TransactionID {
			return &trans, i, true
		}
	}
	return nil, -1, false
}

// deleteAllTransactionsByUid deletes all transactions
// of a user given a user ID.
func (db *Database) deleteAllTransactionsByUid(uid int) {
	transactions := []Transaction{}
	db.Mu.Lock()
	for _, trans := range db.Transactions {
		if uid != trans.UserID {
			transactions = append(transactions, trans)
		}
	}
	db.Transactions = transactions
	db.Mu.Unlock()
}

// findCategoryByCid returns a Category pointer given
// a category ID. It is used for checking existing categories. If
// a category doesn't exist, it returns a nil pointer.
func (db *Database) findCategoryByCid(cid int) *Category {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	for _, cat := range db.Categories {
		if cid == cat.CategoryID {
			return &cat
		}
	}
	return nil
}

// validateNewTransaction validates a POST or PUT transaction request.
// It sends a message to the client if it is a bad request.
func (db *Database) validateNewTransaction(w http.ResponseWriter, r *http.Request, trans Transaction) bool {
	if tempCategory := db.findCategoryByCid(trans.CategoryID); tempCategory == nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, "Category doesn't exist.")
		w.WriteHeader(http.StatusNotFound)
		utils.SendMessageWithBody(w, false, "Category doesn't exist.")
		return false
	}
	if _, err := time.Parse("01-02-2006", trans.Date); err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		utils.SendMessageWithBody(w, false, "Invalid date format.")
		return false
	}
	return true
}
