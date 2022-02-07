package models

import (
	"encoding/json"
	"final-project/server/utils"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Database struct {
	Users             []User        `json:"users"`
	Credentials       []Credentials `json:"credentials"`
	Transactions      []Transaction `json:"transactions"`
	Categories        []Category    `json:"categories"`
	NextUserID        int           `json:"nextUserID"`
	NextTransactionID int           `json:"nextTransactionID"`
	CurrentUserID     int
	Mu                sync.Mutex
}

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
				utils.SendMessageWithBody(w, true, "Logged in successfully!")
				db.CurrentUserID = userID
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
