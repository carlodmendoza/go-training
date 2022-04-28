package transactions

import (
	"encoding/json"
	"strconv"

	gohttp "net/http"

	"github.com/carlodmendoza/go-training/final-project/server/internal/auth"
	"github.com/carlodmendoza/go-training/final-project/server/pkg/http"
	"github.com/carlodmendoza/go-training/final-project/server/storage"

	"github.com/go-chi/chi/v5"
)

func UpdateHandler(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) error {
	username := auth.GetUser(r)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	_, ok, err := db.FindTransaction(username, id)
	if !ok {
		return http.StatusError{Code: gohttp.StatusNotFound, Err: ErrTransactionNotFound}
	}
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	var transactionReq TransactionRequest

	err = json.NewDecoder(r.Body).Decode(&transactionReq)
	if err != nil {
		return http.StatusError{Code: gohttp.StatusBadRequest, Err: err}
	}

	err = validateHandler(db, rw, r, transactionReq)
	if err != nil {
		return err
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
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	_, _ = rw.WriteMessage("Transaction updated successfully!")

	return nil
}
