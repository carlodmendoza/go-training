package transactions

import (
	"strconv"

	gohttp "net/http"

	"github.com/carlodmendoza/go-training/final-project/server/internal/auth"
	"github.com/carlodmendoza/go-training/final-project/server/pkg/http"
	"github.com/carlodmendoza/go-training/final-project/server/storage"

	"github.com/go-chi/chi/v5"
)

func DeleteHandler(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) error {
	username := auth.GetUser(r)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	_, ok, err := db.FindTransaction(username, id)
	if !ok {
		return http.StatusError{Code: gohttp.StatusNotFound, Err: ErrTransactionNotFound}
	}
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	err = db.DeleteTransaction(username, id)
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	_, _ = rw.WriteMessage("Transaction deleted successfully.")

	return nil
}

func Purge(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) error {
	username := auth.GetUser(r)

	ok, err := db.DeleteTransactions(username)
	if !ok {
		return http.StatusError{Code: gohttp.StatusNotFound, Err: ErrTransactionNotFound}
	}
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	_, _ = rw.WriteMessage("All transactions deleted successfully.")

	return nil
}
