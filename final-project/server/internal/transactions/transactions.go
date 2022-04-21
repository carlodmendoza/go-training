package transactions

import (
	"errors"
	"fmt"
	"net/http"
	"server/internal/categories"
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

// validateTransaction validates a POST or PUT transaction request.
// It sends a message to the client if it is a bad request.
func validateTransactionRequest(db storage.Service, w http.ResponseWriter, r *http.Request, transReq TransactionRequest) bool {
	if transReq.Amount == 0 || transReq.Date == "" || transReq.CategoryID == 0 {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
		http.Error(w, "Missing required fields.", http.StatusBadRequest)
		return false
	}

	exists, err := db.CategoryExists(transReq.CategoryID)
	if !exists {
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
