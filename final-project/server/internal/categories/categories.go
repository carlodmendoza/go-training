package categories

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/carlodmendoza/go-training/final-project/server/storage"
)

var (
	ErrInvalidCategory = errors.New("Category doesn't exist")
)

// ProcessCategories handles a categories/ request by a client.
// The client can get all categories.
func ProcessCategories(db storage.Service, w http.ResponseWriter, r *http.Request) {
	categories, err := db.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(categories)
	if err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
