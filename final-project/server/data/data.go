package data

import "final-project/server/models"

var Transactions = []models.Transaction{
	{
		TransactionID: 1,
		Amount:        100,
		Date:          "02-10-2022",
		Notes:         "",
		CategoryID:    1,
	},
	{
		TransactionID: 2,
		Amount:        200,
		Date:          "02-10-2022",
		Notes:         "",
		CategoryID:    2,
	},
	{
		TransactionID: 3,
		Amount:        300,
		Date:          "02-10-2022",
		Notes:         "",
		CategoryID:    3,
	},
	{
		TransactionID: 4,
		Amount:        400,
		Date:          "02-10-2022",
		Notes:         "",
		CategoryID:    4,
	},
	{
		TransactionID: 5,
		Amount:        500,
		Date:          "02-10-2022",
		Notes:         "",
		CategoryID:    5,
	},
	{
		TransactionID: 6,
		Amount:        200,
		Date:          "02-10-2022",
		Notes:         "",
		CategoryID:    17,
	},
	{
		TransactionID: 7,
		Amount:        400,
		Date:          "02-10-2022",
		Notes:         "",
		CategoryID:    18,
	},
	{
		TransactionID: 8,
		Amount:        500,
		Date:          "02-10-2022",
		Notes:         "",
		CategoryID:    19,
	},
	{
		TransactionID: 9,
		Amount:        800,
		Date:          "02-10-2022",
		Notes:         "",
		CategoryID:    20,
	},
	{
		TransactionID: 10,
		Amount:        1000,
		Date:          "02-10-2022",
		Notes:         "",
		CategoryID:    21,
	},
}

var Categories = []models.Category{
	{
		CategoryID: 1,
		Name:       "Food and Beverages",
		Type:       "Expense",
	},
	{
		CategoryID: 2,
		Name:       "Gifts and Donations",
		Type:       "Expense",
	},
	{
		CategoryID: 3,
		Name:       "Education",
		Type:       "Expense",
	},
	{
		CategoryID: 4,
		Name:       "Fees and Charges",
		Type:       "Expense",
	},
	{
		CategoryID: 5,
		Name:       "Bills and Utilities",
		Type:       "Expense",
	},
	{
		CategoryID: 6,
		Name:       "Transportation",
		Type:       "Expense",
	},
	{
		CategoryID: 7,
		Name:       "Shopping",
		Type:       "Expense",
	},
	{
		CategoryID: 8,
		Name:       "Friends and Lover",
		Type:       "Expense",
	},
	{
		CategoryID: 9,
		Name:       "Entertainment",
		Type:       "Expense",
	},
	{
		CategoryID: 10,
		Name:       "Travel",
		Type:       "Expense",
	},
	{
		CategoryID: 11,
		Name:       "Health and Fitness",
		Type:       "Expense",
	},
	{
		CategoryID: 12,
		Name:       "Family",
		Type:       "Expense",
	},
	{
		CategoryID: 13,
		Name:       "Investment",
		Type:       "Expense",
	},
	{
		CategoryID: 14,
		Name:       "Business",
		Type:       "Expense",
	},
	{
		CategoryID: 15,
		Name:       "Insurances",
		Type:       "Expense",
	},
	{
		CategoryID: 16,
		Name:       "Other Expense",
		Type:       "Expense",
	},
	{
		CategoryID: 17,
		Name:       "Salary",
		Type:       "Income",
	},
	{
		CategoryID: 18,
		Name:       "Award",
		Type:       "Income",
	},
	{
		CategoryID: 19,
		Name:       "Interest Money",
		Type:       "Income",
	},
	{
		CategoryID: 20,
		Name:       "Gifts",
		Type:       "Income",
	},
	{
		CategoryID: 21,
		Name:       "Selling",
		Type:       "Income",
	},
	{
		CategoryID: 22,
		Name:       "Other Income",
		Type:       "Income",
	},
}
