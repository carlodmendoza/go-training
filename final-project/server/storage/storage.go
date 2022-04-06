package storage

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
	UserID    int    `json:"uid"`
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
	UserID     int     `json:"uid"`
	CategoryID int     `json:"cid"`
}

type StorageService interface {
	// CreateUser creates a new User with ID, name, and password.
	CreateUser(username, password string)

	// FindUser returns true if a given username already has an existing account.
	// Otherwise, it returns false.
	FindUser(username string) bool

	// AuthenticateUser returns the User ID and true if given username and password is correct.
	// Otherwise, it returns 0 and false.
	AuthenticateUser(username, password string) (int, bool)

	// CreateSession creates a new Session or updates an existing Session given the User ID,
	// and returns the Session.
	CreateSession(uid int) Session

	// FindSession returns the associated User ID given a Session token.
	// If no associated User is found, it returns 0.
	FindSession(token string) int

	// GetCategories returns the list of Category.
	GetCategories() []Category

	// FindCategory returns true if a given Category ID exists.
	// Otherwise, it returns false.
	FindCategory(cid int) bool

	// CreateTransaction creates a new Transaction and associates it to a User ID.
	CreateTransaction(tr Transaction)

	// GetTransactions returns a list of Transaction given a User ID.
	// If there are no existing user transactions, it returns an empty list.
	GetTransactions(uid int) []Transaction

	// UpdateTransaction updates an existing Transaction.
	// The Transaction ID and associated User ID should not change.
	UpdateTransaction(tr Transaction)

	// DeleteTransactions deletes all transactions of a User given the User ID.
	// If there are no transactions to delete, it returns false; otherwise, it returns true.
	DeleteTransactions(uid int) bool

	// DeleteTransaction deletes a Transaction of a User given the User ID and Transaction ID.
	DeleteTransaction(uid, tid int)

	// FindTransaction returns a Transaction and true if given Transaction ID exists in given user ID.
	// Otherwise, it returns an empty Transaction and false.
	FindTransaction(uid, tid int) (Transaction, bool)
}
