package transactions

import (
	"net/http"
	"server/internal/auth"
	"server/storage"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func DeleteHandler(db storage.Service, w http.ResponseWriter, r *http.Request) {
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

	err = db.DeleteTransaction(username, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("Transaction deleted successfully."))
}

func Purge(db storage.Service, w http.ResponseWriter, r *http.Request) {
	username := auth.GetUser(r)

	ok, err := db.DeleteTransactions(username)
	if !ok {
		http.Error(w, ErrTransactionNotFound.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("All transactions deleted successfully."))
}
