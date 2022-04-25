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

func RetrieveHandler(db storage.Service, w http.ResponseWriter, r *http.Request) {
	username := auth.GetUser(r)
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	transaction, ok, err := db.FindTransaction(username, id)
	if !ok {
		http.Error(w, ErrTransactionNotFound.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(transaction)
	if err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
