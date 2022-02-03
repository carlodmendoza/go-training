package models

type User struct {
	UserID            int
	Name              string
	Credentials       Credentials
	Transactions      []Transaction
	NextTransactionID int
}

type Credentials struct {
	Username string
	Password string
}
