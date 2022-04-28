package categories

import (
	"encoding/json"
	"errors"
	gohttp "net/http"

	"github.com/carlodmendoza/go-training/final-project/server/pkg/http"
	"github.com/carlodmendoza/go-training/final-project/server/storage"
)

var (
	ErrInvalidCategory = errors.New("Category doesn't exist")
)

// ProcessCategories handles a categories/ request by a client.
// The client can get all categories.
func ProcessCategories(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) error {
	categories, err := db.GetCategories()
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	out, err := json.Marshal(categories)
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	_, _ = rw.Write(out)

	return nil
}
