package models

type User struct {
	UserID int    `json:"userID"`
	Name   string `json:"name"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserID   int    `json:"userID"`
}

type Transaction struct {
	TransactionID int     `json:"transactionID"`
	Amount        float64 `json:"amount"`
	Date          string  `json:"date"`
	Notes         string  `json:"notes"`
	UserID        int     `json:"userID"`
	CategoryID    int     `json:"categoryID"`
}

type Category struct {
	CategoryID int    `json:"categoryID"`
	Name       string `json:"name"`
	Type       string `json:"type"`
}
