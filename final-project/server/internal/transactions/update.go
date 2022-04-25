package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/carlodmendoza/go-training/final-project/server/internal/auth"
	"github.com/carlodmendoza/go-training/final-project/server/storage"

	"github.com/go-chi/chi/v5"
)

func UpdateHandler(db storage.Service, w http.ResponseWriter, r *http.Request) {
	username := auth.GetUser(r)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	_, ok, err := db.FindTransaction(username, id)
	if !ok {
		http.Error(w, ErrTransactionNotFound.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var transactionReq TransactionRequest

	err = json.NewDecoder(r.Body).Decode(&transactionReq)
	if err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !validateTransactionRequest(db, w, r, transactionReq) {
		return
	}

	transaction := storage.Transaction{
		ID:         id,
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
}
