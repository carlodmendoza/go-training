package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/carlodmendoza/go-training/final-project/server/internal/auth"
	"github.com/carlodmendoza/go-training/final-project/server/storage"
)

func CreateHandler(db storage.Service, w http.ResponseWriter, r *http.Request) {
	username := auth.GetUser(r)

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
}
