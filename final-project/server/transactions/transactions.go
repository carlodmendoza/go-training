package transactions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/categories"
	"server/storage"
	"time"
)

type TransactionRequest struct {
	Amount     float64 `json:"amount"`
	Date       string  `json:"date"`
	Notes      string  `json:"notes,omitempty"`
	CategoryID int     `json:"category_id"`
}

var (
	ErrTransactionNotFound = errors.New("No transaction/s found")
)

// ProcessTransaction handles a transaction/ request by a client
// given a username. The client can either get all transactions,
// add new transaction, or delete all transactions.
func ProcessTransaction(db storage.StorageService, w http.ResponseWriter, r *http.Request, username string) {
	switch r.Method {
	case "GET":
		transactions, err := db.GetTransactions(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(transactions)
		if err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "POST":
		var transactionReq TransactionRequest

		err := json.NewDecoder(r.Body).Decode(&transactionReq)
		if err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !validateTransactionRequest(db, w, r, transactionReq) {
			return
		}

		transaction := storage.Transaction{
			Amount:     transactionReq.Amount,
			Date:       transactionReq.Date,
			Notes:      transactionReq.Notes,
			Username:   username,
			CategoryID: transactionReq.CategoryID,
		}
		err = db.CreateTransaction(transaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("Transaction added successfully!"))
	case "DELETE":
		ok, err := db.DeleteTransactions(username)
		if !ok {
			http.Error(w, ErrTransactionNotFound.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write([]byte("All transactions deleted successfully."))
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// ProcessTransactionID handles a transaction/id request by a client
// given a username and a transaction ID. The client can either get,
// update, or delete a transaction.
func ProcessTransactionID(db storage.StorageService, w http.ResponseWriter, r *http.Request, username string, transID int) {
	transaction, ok, err := db.FindTransaction(username, transID)
	if !ok {
		http.Error(w, ErrTransactionNotFound.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		err := json.NewEncoder(w).Encode(transaction)
		if err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "PUT":
		var transactionReq TransactionRequest

		err := json.NewDecoder(r.Body).Decode(&transactionReq)
		if err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !validateTransactionRequest(db, w, r, transactionReq) {
			return
		}

		transaction := storage.Transaction{
			ID:         transID,
			Amount:     transactionReq.Amount,
			Date:       transactionReq.Date,
			Notes:      transactionReq.Notes,
			Username:   username,
			CategoryID: transactionReq.CategoryID,
		}
		err = db.UpdateTransaction(transaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write([]byte("Transaction updated successfully!"))
	case "DELETE":
		err := db.DeleteTransaction(username, transID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write([]byte("Transaction deleted successfully."))
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// validateTransaction validates a POST or PUT transaction request.
// It sends a message to the client if it is a bad request.
func validateTransactionRequest(db storage.StorageService, w http.ResponseWriter, r *http.Request, transReq TransactionRequest) bool {
	if transReq.Amount == 0 || transReq.Date == "" || transReq.CategoryID == 0 {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
		http.Error(w, "Missing required fields.", http.StatusBadRequest)
		return false
	}

	ok, err := db.FindCategory(transReq.CategoryID)
	if !ok {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, categories.ErrInvalidCategory)
		http.Error(w, categories.ErrInvalidCategory.Error(), http.StatusNotFound)
		return false
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	if _, err := time.Parse("01-02-2006", transReq.Date); err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return false
	}

	return true
}
