package transactions

import (
	"errors"
	gohttp "net/http"
	"time"

	"github.com/carlodmendoza/go-training/final-project/server/internal/categories"
	"github.com/carlodmendoza/go-training/final-project/server/pkg/http"
	"github.com/carlodmendoza/go-training/final-project/server/storage"
)

type TransactionRequest struct {
	Amount     float64 `json:"amount"`
	Date       string  `json:"date"`
	Notes      string  `json:"notes,omitempty"`
	CategoryID int     `json:"category_id"`
}

var (
	ErrTransactionNotFound = errors.New("No transaction/s found")
	ErrEmptyFields         = errors.New("Transaction amount, date, or category ID is empty")
)

// validateHandler validates a POST or PUT transaction request.
// It sends a message to the client if it is a bad request.
func validateHandler(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request, transReq TransactionRequest) (int, error) {
	if transReq.Amount == 0 || transReq.Date == "" || transReq.CategoryID == 0 {
		return gohttp.StatusBadRequest, ErrEmptyFields
	}

	exists, err := db.CategoryExists(transReq.CategoryID)
	if !exists {
		return gohttp.StatusNotFound, categories.ErrInvalidCategory
	}
	if err != nil {
		return gohttp.StatusInternalServerError, err
	}

	_, err = time.Parse("01-02-2006", transReq.Date)
	if err != nil {
		return gohttp.StatusBadRequest, err
	}

	return gohttp.StatusContinue, nil
}
