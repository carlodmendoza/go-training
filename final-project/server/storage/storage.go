package storage

type User struct {
	ID           int              `json:"id"`
	Name         string           `json:"username"`
	Password     string           `json:"password"`
	SessionToken string           `json:"session_token"`
	Transactions map[int]struct{} `json:"transaction_ids"`
}

type Session struct {
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
	Username  string `json:"username"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Transaction struct {
	ID         int     `json:"id"`
	Amount     float64 `json:"amount"`
	Date       string  `json:"date"`
	Notes      string  `json:"notes"`
	Username   string  `json:"username"`
	CategoryID int     `json:"category_id"`
}

type Service interface {
	// CreateUser creates a new User with ID, name, password, empty session token, and empty transaction IDs.
	CreateUser(username, password string) error

	// FindUser returns true if a given username already has an existing account.
	// Otherwise, it returns false.
	FindUser(username string) (bool, error)

	// AuthenticateUser returns true if given username and password is correct.
	// Otherwise, it returns false.
	AuthenticateUser(username, password string) (bool, error)

	// CreateSession creates a new Session or updates an existing Session given the username and generated token.
	// It also associates a Session to a User.
	CreateSession(username, token string) error

	// FindSession returns the associated username given a Session token.
	// If no associated User is found, it returns an empty string.
	FindSession(token string) (string, error)

	// GetCategories returns the list of Category.
	GetCategories() ([]Category, error)

	// FindCategory returns true if a given Category ID exists.
	// Otherwise, it returns false.
	FindCategory(cid int) (bool, error)

	// CreateTransaction creates a new Transaction and associates it to a User.
	CreateTransaction(tr Transaction) error

	// GetTransactions returns a list of Transaction given a username.
	// If there are no existing user transactions, it returns an empty list.
	GetTransactions(username string) ([]Transaction, error)

	// UpdateTransaction updates an existing Transaction.
	// The Transaction ID and associated username should not change.
	UpdateTransaction(tr Transaction) error

	// DeleteTransactions deletes all transactions of a User given the username.
	// If there are no transactions to delete, it returns false; otherwise, it returns true.
	DeleteTransactions(username string) (bool, error)

	// DeleteTransaction deletes a Transaction of a User given the username and Transaction ID.
	DeleteTransaction(username string, tid int) error

	// FindTransaction returns a Transaction and true if given Transaction ID exists in given username.
	// Otherwise, it returns an empty Transaction and false.
	FindTransaction(username string, tid int) (Transaction, bool, error)
}
