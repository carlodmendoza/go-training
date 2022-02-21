package models

/*
	Models program contains all structs being
	used by the server.
	Author: Carlo Mendoza
*/

type User struct {
	UserID int `json:"userID"`
}

type Session struct {
	Token     string `json:"token"`
	Timestamp int64  `json:"timestamp"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Transaction struct {
	Amount     float64 `json:"amount"`
	Date       string  `json:"date"`
	Notes      string  `json:"notes"`
	CategoryID int     `json:"categoryID"`
}

type Category struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
