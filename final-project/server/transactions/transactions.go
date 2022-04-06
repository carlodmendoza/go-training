package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/storage"
	"time"
)

type TransactionRequest struct {
	Amount     float64 `json:"amount"`
	Date       string  `json:"date"`
	Notes      string  `json:"notes,omitempty"`
	CategoryID int     `json:"category_id"`
}

// ProcessTransaction handles a transaction/ request by a client
// given a user ID. The client can either get all transactions,
// add new transaction, or delete all transactions.
func ProcessTransaction(db storage.StorageService, w http.ResponseWriter, r *http.Request, userID int) {
	switch r.Method {
	case "GET":
		err := json.NewEncoder(w).Encode(db.GetTransactions(userID))
		if err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "POST":
		var transactionReq TransactionRequest

		err := json.NewDecoder(r.Body).Decode(&transactionReq)
		if err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if transactionReq.Amount == 0 || transactionReq.Date == "" || transactionReq.CategoryID == 0 {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			http.Error(w, "Missing required fields.", http.StatusBadRequest)
			return
		}

		if !validateTransactionRequest(db, w, r, transactionReq) {
			return
		}

		transaction := storage.Transaction{
			Amount:     transactionReq.Amount,
			Date:       transactionReq.Date,
			Notes:      transactionReq.Notes,
			UserID:     userID,
			CategoryID: transactionReq.CategoryID,
		}
		db.CreateTransaction(transaction)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("Transaction added successfully!"))
	case "DELETE":
		if !db.DeleteTransactions(userID) {
			http.Error(w, "No transactions found.", http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte("All transactions deleted successfully."))
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// ProcessTransactionID handles a transaction/id request by a client
// given a user ID and a transaction ID. The client can either get,
// update, or delete a transaction.
func ProcessTransactionID(db storage.StorageService, w http.ResponseWriter, r *http.Request, userID, transID int) {
	transaction, ok := db.FindTransaction(userID, transID)
	if !ok {
		http.Error(w, "Transaction not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		err := json.NewEncoder(w).Encode(transaction)
		if err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "PUT":
		var transactionReq TransactionRequest

		err := json.NewDecoder(r.Body).Decode(&transactionReq)
		if err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
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
			UserID:     userID,
			CategoryID: transactionReq.CategoryID,
		}
		db.UpdateTransaction(transaction)
		_, _ = w.Write([]byte("Transaction updated successfully!"))
	case "DELETE":
		db.DeleteTransaction(userID, transID)
		_, _ = w.Write([]byte("Transaction deleted successfully."))
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// validateTransaction validates a POST or PUT transaction request.
// It sends a message to the client if it is a bad request.
func validateTransactionRequest(db storage.StorageService, w http.ResponseWriter, r *http.Request, transReq TransactionRequest) bool {
	if !db.FindCategory(transReq.CategoryID) {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, "Category doesn't exist.")
		http.Error(w, "Category doesn't exist.", http.StatusNotFound)
		return false
	}
	if _, err := time.Parse("01-02-2006", transReq.Date); err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}
