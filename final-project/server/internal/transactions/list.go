package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/carlodmendoza/go-training/final-project/server/internal/auth"
	"github.com/carlodmendoza/go-training/final-project/server/storage"
)

func ListHandler(db storage.Service, w http.ResponseWriter, r *http.Request) {
	username := auth.GetUser(r)

	transactions, err := db.GetTransactions(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
