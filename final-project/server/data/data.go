package data

import "final-project/server/models"

var Transactions = []models.Transaction{
	{
		Amount:     100,
		Date:       "02-10-2022",
		Notes:      "",
		CategoryID: 1,
	},
	{
		Amount:     200,
		Date:       "02-10-2022",
		Notes:      "",
		CategoryID: 2,
	},
	{
		Amount:     300,
		Date:       "02-10-2022",
		Notes:      "",
		CategoryID: 3,
	},
	{
		Amount:     400,
		Date:       "02-10-2022",
		Notes:      "",
		CategoryID: 4,
	},
	{
		Amount:     500,
		Date:       "02-10-2022",
		Notes:      "",
		CategoryID: 5,
	},
	{
		Amount:     200,
		Date:       "02-10-2022",
		Notes:      "",
		CategoryID: 17,
	},
	{
		Amount:     400,
		Date:       "02-10-2022",
		Notes:      "",
		CategoryID: 18,
	},
	{
		Amount:     500,
		Date:       "02-10-2022",
		Notes:      "",
		CategoryID: 19,
	},
	{
		Amount:     800,
		Date:       "02-10-2022",
		Notes:      "",
		CategoryID: 20,
	},
	{
		Amount:     1000,
		Date:       "02-10-2022",
		Notes:      "",
		CategoryID: 21,
	},
}

var Categories = []models.Category{
	{
		Name: "Food and Beverages",
		Type: "Expense",
	},
	{
		Name: "Gifts and Donations",
		Type: "Expense",
	},
	{
		Name: "Education",
		Type: "Expense",
	},
	{
		Name: "Fees and Charges",
		Type: "Expense",
	},
	{
		Name: "Bills and Utilities",
		Type: "Expense",
	},
	{
		Name: "Transportation",
		Type: "Expense",
	},
	{
		Name: "Shopping",
		Type: "Expense",
	},
	{
		Name: "Friends and Lover",
		Type: "Expense",
	},
	{
		Name: "Entertainment",
		Type: "Expense",
	},
	{
		Name: "Travel",
		Type: "Expense",
	},
	{
		Name: "Health and Fitness",
		Type: "Expense",
	},
	{
		Name: "Family",
		Type: "Expense",
	},
	{
		Name: "Investment",
		Type: "Expense",
	},
	{
		Name: "Business",
		Type: "Expense",
	},
	{
		Name: "Insurances",
		Type: "Expense",
	},
	{
		Name: "Other Expense",
		Type: "Expense",
	},
	{
		Name: "Salary",
		Type: "Income",
	},
	{
		Name: "Award",
		Type: "Income",
	},
	{
		Name: "Interest Money",
		Type: "Income",
	},
	{
		Name: "Gifts",
		Type: "Income",
	},
	{
		Name: "Selling",
		Type: "Income",
	},
	{
		Name: "Other Income",
		Type: "Income",
	},
}
