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

func RetrieveHandler(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) error {
	username := auth.GetUser(r)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	transaction, ok, err := db.FindTransaction(username, id)
	if !ok {
		return http.StatusError{Code: gohttp.StatusNotFound, Err: ErrTransactionNotFound}
	}
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	out, err := json.Marshal(transaction)
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	_, _ = rw.Write(out)

	return nil
}
