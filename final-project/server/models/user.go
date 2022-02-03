package models

type User struct {
	UserID            int           `json:"userID"`
	Name              string        `json:"name"`
	Credentials       Credentials   `json:"credentials"`
	Transactions      []Transaction `json:"transactions"`
	NextTransactionID int           `json:"nextTransactionID"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Transaction struct {
	TransactionID string  `json:"transactionID"`
	Type          string  `json:"type"`
	Category      string  `json:"category"`
	Amount        float64 `json:"amount"`
	Date          string  `json:"date"`
	Notes         string  `json:"notes"`
}
