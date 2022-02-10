package models

/*
	Models program contains all structs being
	used by the server.
	Author: Carlo Mendoza
*/

type Response struct {
	Message string
	Success bool
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
