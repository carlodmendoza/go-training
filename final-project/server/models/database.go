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
			if tempCategory := db.findCategoryByCid(transaction.CategoryID); tempCategory == nil {
				fmt.Printf("Error in %s: %s\n", r.URL.Path, "Category doesn't exist.")
				w.WriteHeader(http.StatusBadRequest)
				utils.SendMessageWithBody(w, false, "Category doesn't exist.")
				return
			}
			if _, err := time.Parse("01-02-2006", transaction.Date); err != nil {
				fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
				w.WriteHeader(http.StatusBadRequest)
				utils.SendMessageWithBody(w, false, "Invalid date format.")
				return
			}

			db.Mu.Lock()
			db.NextTransactionID++
			transaction.TransactionID = db.NextTransactionID
			transaction.UserID = userID
			db.Transactions = append(db.Transactions, transaction)
			db.Mu.Unlock()

			w.WriteHeader(http.StatusCreated)
			utils.SendMessageWithBody(w, true, "Transaction added successfully!")
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

func (db *Database) ProcessTransactionID(w http.ResponseWriter, r *http.Request, id, transID int) {
	switch r.Method {
	case "GET":
	case "PUT":
	case "DELETE":
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

func (db *Database) findCategoryByCid(catID int) *Category {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	for _, cat := range db.Categories {
		if catID == cat.CategoryID {
			return &cat
		}
	}
	return nil
}

// func (user *User) retrieveTransactionById(id int) (*Transaction, int, bool) {
// 	for i, tran := range user.Transactions {
// 		if tran.TransactionID == id {
// 			return &tran, i, true
// 		}
// 	}
// 	return nil, -1, false
// }
