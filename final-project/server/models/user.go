package models

import (
	"encoding/json"
	"final-project/server/utils"
	"fmt"
	"net/http"
)

type User struct {
	UserID       int           `json:"userID"`
	Name         string        `json:"name"`
	Credentials  Credentials   `json:"credentials"`
	Transactions []Transaction `json:"transactions"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Transaction struct {
	TransactionID int     `json:"transactionID"`
	Type          string  `json:"type"`
	Category      string  `json:"category"`
	Amount        float64 `json:"amount"`
	Date          string  `json:"date"`
	Notes         string  `json:"notes"`
}

func (user *User) ProcessTransaction(db *Database, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		db.Mu.Lock()
		if err := json.NewEncoder(w).Encode(user.Transactions); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			utils.SendMessageWithBody(w, false, "500 Internal Server Error")
			return
		}
		db.Mu.Unlock()
	case "POST":
		var transaction Transaction
		if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		if transaction.Type == "" || transaction.Category == "" || transaction.Amount == 0 || transaction.Date == "" {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
		} else {
			db.Mu.Lock()
			db.NextTransactionID++
			transaction.TransactionID = db.NextTransactionID
			user.Transactions = append(user.Transactions, transaction)
			db.Mu.Unlock()

			w.WriteHeader(http.StatusCreated)
			utils.SendMessageWithBody(w, true, "Transaction added successfully!")
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

func (user *User) ProcessTransactionID(id int, db *Database, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	case "PUT":
	case "DELETE":
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

func (user *User) retrieveTransactionById(id int) (*Transaction, int, bool) {
	for i, tran := range user.Transactions {
		if tran.TransactionID == id {
			return &tran, i, true
		}
	}
	return nil, -1, false
}
