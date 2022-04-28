package transactions

import (
	"encoding/json"
	gohttp "net/http"

	"github.com/carlodmendoza/go-training/final-project/server/internal/auth"
	"github.com/carlodmendoza/go-training/final-project/server/pkg/http"
	"github.com/carlodmendoza/go-training/final-project/server/storage"
)

func CreateHandler(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) error {
	username := auth.GetUser(r)

	var transactionReq TransactionRequest

	err := json.NewDecoder(r.Body).Decode(&transactionReq)
	if err != nil {
		return http.StatusError{Code: gohttp.StatusBadRequest, Err: err}
	}

	err = validateHandler(db, rw, r, transactionReq)
	if err != nil {
		return err
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
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	rw.WriteHeader(gohttp.StatusCreated)
	_, _ = rw.WriteMessage("Transaction added successfully!")

	return nil
}
