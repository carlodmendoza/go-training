package transactions

import (
	"encoding/json"
	gohttp "net/http"

	"github.com/carlodmendoza/go-training/final-project/server/internal/auth"
	"github.com/carlodmendoza/go-training/final-project/server/pkg/http"
	"github.com/carlodmendoza/go-training/final-project/server/storage"
)

func ListHandler(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) error {
	username := auth.GetUser(r)

	transactions, err := db.GetTransactions(username)
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	out, err := json.Marshal(transactions)
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	_, _ = rw.Write(out)

	return nil
}
