package models

type Transaction struct {
	TransactionID string
	Category      Category
	Amount        float64
	Date          string
	Notes         string
}

type Category struct {
	Name string
	Type string
}
