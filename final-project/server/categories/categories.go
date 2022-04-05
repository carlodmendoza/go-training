package categories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/storage"
)

// ProcessCategories handles a categories/ request by a client.
// The client can get all categories.
func ProcessCategories(db storage.StorageService, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := json.NewEncoder(w).Encode(db.GetCategories())
		if err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
