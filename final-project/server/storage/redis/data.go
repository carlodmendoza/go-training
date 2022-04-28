package redis

import (
	"github.com/carlodmendoza/go-training/final-project/server/storage"
)

var Transactions = []storage.Transaction{
	{
		ID:         1,
		Amount:     100,
		Date:       "02-10-2022",
		Notes:      "",
		Username:   "cmendoza",
		CategoryID: 1,
	},
	{
		ID:         2,
		Amount:     200,
		Date:       "02-10-2022",
		Notes:      "",
		Username:   "cmendoza",
		CategoryID: 2,
	},
	{
		ID:         3,
		Amount:     300,
		Date:       "02-10-2022",
		Notes:      "",
		Username:   "cmendoza",
		CategoryID: 3,
	},
	{
		ID:         4,
		Amount:     400,
		Date:       "02-10-2022",
		Notes:      "",
		Username:   "cmendoza",
		CategoryID: 4,
	},
	{
		ID:         5,
		Amount:     500,
		Date:       "02-10-2022",
		Notes:      "",
		Username:   "cmendoza",
		CategoryID: 5,
	},
	{
		ID:         6,
		Amount:     200,
		Date:       "02-10-2022",
		Notes:      "",
		Username:   "cmendoza",
		CategoryID: 17,
	},
	{
		ID:         7,
		Amount:     400,
		Date:       "02-10-2022",
		Notes:      "",
		Username:   "cmendoza",
		CategoryID: 18,
	},
	{
		ID:         8,
		Amount:     500,
		Date:       "02-10-2022",
		Notes:      "",
		Username:   "cmendoza",
		CategoryID: 19,
	},
	{
		ID:         9,
		Amount:     800,
		Date:       "02-10-2022",
		Notes:      "",
		Username:   "cmendoza",
		CategoryID: 20,
	},
	{
		ID:         10,
		Amount:     1000,
		Date:       "02-10-2022",
		Notes:      "",
		Username:   "cmendoza",
		CategoryID: 21,
	},
}

var Categories = []storage.Category{
	{
		ID:   1,
		Name: "Food and Beverages",
		Type: "Expense",
	},
	{
		ID:   2,
		Name: "Gifts and Donations",
		Type: "Expense",
	},
	{
		ID:   3,
		Name: "Education",
		Type: "Expense",
	},
	{
		ID:   4,
		Name: "Fees and Charges",
		Type: "Expense",
	},
	{
		ID:   5,
		Name: "Bills and Utilities",
		Type: "Expense",
	},
	{
		ID:   6,
		Name: "Transportation",
		Type: "Expense",
	},
	{
		ID:   7,
		Name: "Shopping",
		Type: "Expense",
	},
	{
		ID:   8,
		Name: "Friends and Lover",
		Type: "Expense",
	},
	{
		ID:   9,
		Name: "Entertainment",
		Type: "Expense",
	},
	{
		ID:   10,
		Name: "Travel",
		Type: "Expense",
	},
	{
		ID:   11,
		Name: "Health and Fitness",
		Type: "Expense",
	},
	{
		ID:   12,
		Name: "Family",
		Type: "Expense",
	},
	{
		ID:   13,
		Name: "Investment",
		Type: "Expense",
	},
	{
		ID:   14,
		Name: "Business",
		Type: "Expense",
	},
	{
		ID:   15,
		Name: "Insurances",
		Type: "Expense",
	},
	{
		ID:   16,
		Name: "Other Expense",
		Type: "Expense",
	},
	{
		ID:   17,
		Name: "Salary",
		Type: "Income",
	},
	{
		ID:   18,
		Name: "Award",
		Type: "Income",
	},
	{
		ID:   19,
		Name: "Interest Money",
		Type: "Income",
	},
	{
		ID:   20,
		Name: "Gifts",
		Type: "Income",
	},
	{
		ID:   21,
		Name: "Selling",
		Type: "Income",
	},
	{
		ID:   22,
		Name: "Other Income",
		Type: "Income",
	},
}
