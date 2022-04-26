package transactions

import (
	"encoding/json"
	gohttp "net/http"

	"github.com/carlodmendoza/go-training/final-project/server/internal/auth"
	"github.com/carlodmendoza/go-training/final-project/server/pkg/http"
	"github.com/carlodmendoza/go-training/final-project/server/storage"
)

func CreateHandler(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) (int, error) {
	username := auth.GetUser(r)

	var transactionReq TransactionRequest

	err := json.NewDecoder(r.Body).Decode(&transactionReq)
	if err != nil {
		return gohttp.StatusBadRequest, err
	}

	status, err := validateHandler(db, rw, r, transactionReq)
	if err != nil {
		return status, err
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
		return gohttp.StatusInternalServerError, err
	}

	_, _ = rw.WriteMessage("Transaction added successfully!")

	return gohttp.StatusCreated, nil
}
