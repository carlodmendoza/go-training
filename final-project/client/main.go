package main

import (
	"bytes"
	"encoding/json"
	"final-project/client/models"
	"final-project/client/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"
)

/*
	Main program for running the client, making requests
	to the server, and generating reports.
	Author: Carlo Mendoza
*/

const baseURL = "http://localhost:8080/"

var cookie *http.Cookie
var categories []models.Category

func main() {
	var choice int
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	utils.PrintWelcomeMessage()
	commands := []string{"Sign in", "Sign up", "Exit"}
	utils.PrintCommands(commands)
	fmt.Scan(&choice)
	for {
		if choice == 1 {
			if signin(c) {
				categories = getCategories(c)
				break
			}
		} else if choice == 2 {
			signup(c)
		} else if choice == 3 {
			utils.PrintExitMessage()
			os.Exit(1)
		}
		utils.PrintCommands(commands)
		fmt.Scan(&choice)
	}

	commands = []string{"View my transactions", "View report", "Add new transaction", "View a transaction", "Edit a transaction", "Delete a transaction", "Delete all transactions", "Exit"}
	utils.PrintCommands(commands)
	fmt.Scan(&choice)
	for {
		if choice == 1 {
			viewTransactions(c, 0)
		} else if choice == 2 {
			generateReport(c)
		} else if choice == 3 {
			addEditTransaction(c, 0)
		} else if choice == 4 {
			var id int
			for id == 0 {
				fmt.Println("Enter the transaction ID: ")
				fmt.Scan(&id)
			}
			viewTransactions(c, id)
		} else if choice == 5 {
			var id int
			for id == 0 {
				fmt.Println("Enter the transaction ID: ")
				fmt.Scan(&id)
			}
			addEditTransaction(c, id)
		} else if choice == 6 {
			var id int
			for id == 0 {
				fmt.Println("Enter the transaction ID: ")
				fmt.Scan(&id)
			}
			deleteTransactions(c, id)
		} else if choice == 7 {
			var confirmChoice int
			confirm := []string{"Go Back", "Continue"}
			for confirmChoice != 1 && confirmChoice != 2 {
				fmt.Println("Warning: This will delete ALL of your transactions. Continue?")
				utils.PrintCommands(confirm)
				fmt.Scan(&confirmChoice)
			}
			if confirmChoice == 2 {
				deleteTransactions(c, 0)
			}
		} else if choice == 8 {
			utils.PrintExitMessage()
			os.Exit(1)
		}
		utils.PrintCommands(commands)
		fmt.Scan(&choice)
	}
}

// signin gets input from user for making a sign in request.
func signin(c http.Client) bool {
	url := baseURL + "signin"

	var username, password string
	fmt.Println("Enter your username: ")
	fmt.Scan(&username)
	fmt.Println("Enter your password: ")
	fmt.Scan(&password)

	reqBody := fmt.Sprintf("{\"username\":\"%s\", \"password\":\"%s\"}", username, password)
	return getResponse(c, url, "POST", reqBody, false)
}

// signup gets input from user for making a sign up request.
func signup(c http.Client) {
	url := baseURL + "signup"

	var name, username, password string
	fmt.Println("Enter your name (no spaces): ")
	fmt.Scan(&name)
	fmt.Println("Enter your username: ")
	fmt.Scan(&username)
	fmt.Println("Enter your password: ")
	fmt.Scan(&password)

	reqBody := fmt.Sprintf("{\"name\":\"%s\", \"username\":\"%s\", \"password\":\"%s\"}", name, username, password)
	getResponse(c, url, "POST", reqBody, false)
}

// viewTransactions prints the user transactions
// received from the server.
func viewTransactions(c http.Client, transID int) {
	var url string
	if transID == 0 {
		url = baseURL + "transactions"
	} else {
		url = baseURL + "transactions/" + fmt.Sprint(transID)
	}
	if trans, ok := getTransactions(c, url, transID); ok {
		if len(trans) > 0 {
			w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
			fmt.Fprintln(w, "ID\tAmount\tDate\tNotes\tCategory")
			for _, t := range trans {
				if t.Notes == "" {
					t.Notes = "\"\""
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", t.TransactionID, utils.FormatFloat(t.Amount), t.Date, t.Notes, getCategoryDetails(t.CategoryID))
			}
			w.Flush()
		} else {
			fmt.Println("No transaction/s found.")
		}
	}
}

// getTransactions handles sending a GET request to get transactions.
// If transaction ID is 0, it gets all transactions; otherwise, it
// gets only a specific transaction.
func getTransactions(c http.Client, url string, transID int) ([]models.Transaction, bool) {
	var transactionList []models.Transaction
	var transaction models.Transaction
	var response models.Response

	req, _ := http.NewRequest("GET", url, nil)
	req.AddCookie(cookie)
	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response: %s", err.Error())
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if transID == 0 {
		if err := json.Unmarshal(data, &transactionList); err != nil {
			fmt.Printf("Failed to parse json response: %s\n", err.Error())
			return []models.Transaction{}, false
		}
		return transactionList, true
	} else {
		if resp.StatusCode == http.StatusNotFound {
			if err := json.Unmarshal(data, &response); err != nil {
				fmt.Printf("Failed to parse json response: %s\n", err.Error())
				return []models.Transaction{}, false
			}
			fmt.Printf("Error: %s\n", response.Message)
			return []models.Transaction{}, false
		}
		if err := json.Unmarshal(data, &transaction); err != nil {
			fmt.Printf("Failed to parse json response: %s\n", err.Error())
			return []models.Transaction{}, false
		}
		return []models.Transaction{transaction}, true
	}
}

// addEditTransaction handles sending a POST or PUT request to add or edit
// a transaction. If transaction ID is 0, it adds a new transaction; otherwise,
// it edits a specific transaction.
func addEditTransaction(c http.Client, transID int) {
	var amount float64
	var url, date, notes string
	var categoryID int

	printCategoryDetails()
	for categoryID == 0 {
		fmt.Println("Enter the category (Enter the ID of your choice): ")
		fmt.Scan(&categoryID)
	}
	for amount == 0 {
		fmt.Println("Enter the amount of transaction: ")
		fmt.Scan(&amount)
	}
	fmt.Println("Enter the date (MM-DD-YYYY) of transaction: ")
	fmt.Scan(&date)
	fmt.Println("Enter any notes (no spaces) about the transaction (Type NA if none): ")
	fmt.Scan(&notes)
	if notes == "NA" {
		notes = ""
	}

	reqBody := fmt.Sprintf("{\"amount\":%f, \"date\":\"%s\", \"notes\":\"%s\", \"categoryID\":%d}", amount, date, notes, categoryID)
	if transID == 0 {
		url = baseURL + "transactions"
		getResponse(c, url, "POST", reqBody, true)
	} else {
		url = baseURL + "transactions/" + fmt.Sprint(transID)
		getResponse(c, url, "PUT", reqBody, true)
	}
}

// deleteTransactions handles sending a DELETE request to delete transactions.
// If transaction ID is 0, it deletes all transactions; otherwise, it
// deletes only a specific transaction.
func deleteTransactions(c http.Client, transID int) {
	var url, reqBody string

	if transID == 0 {
		url = baseURL + "transactions"
	} else {
		url = baseURL + "transactions/" + fmt.Sprint(transID)
	}
	getResponse(c, url, "DELETE", reqBody, true)
}

// getResponse handles sending a general request if client is
// expecting a JSON with message and success fields.
func getResponse(c http.Client, url, method, reqBody string, requireCookie bool) bool {
	var response models.Response
	var req *http.Request

	if method == "POST" || method == "PUT" {
		outData := bytes.NewBuffer([]byte(reqBody))
		req, _ = http.NewRequest(method, url, outData)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	if requireCookie {
		req.AddCookie(cookie)
	}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response: %s", err.Error())
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &response); err != nil {
		fmt.Printf("Failed to parse json response: %s\n", err.Error())
		return false
	}
	if response.Success {
		if url == baseURL+"signin" {
			cookie = resp.Cookies()[0]
		}
		fmt.Println(response.Message)
	} else {
		fmt.Printf("Error: %s\n", response.Message)
	}
	return response.Success
}

// getCategories handles sending a GET request to get categories.
func getCategories(c http.Client) []models.Category {
	url := baseURL + "categories"
	var categories []models.Category

	req, _ := http.NewRequest("GET", url, nil)
	req.AddCookie(cookie)
	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response: %s", err.Error())
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &categories); err != nil {
		fmt.Printf("Failed to parse json response: %s\n", err.Error())
		return []models.Category{}
	}
	return categories
}

// getCategoryDetails returns a formatted string containing
// details of a category given a category ID.
func getCategoryDetails(catID int) string {
	for _, cat := range categories {
		if catID == cat.CategoryID {
			return fmt.Sprintf("%s (%s)", cat.Name, cat.Type)
		}
	}
	return fmt.Sprint(catID)
}

// printCategoryDetails prints a list of all categories as
// reference for user when adding or updating transactions.
func printCategoryDetails() {
	expense := "Expense: \n"
	income := "Income: \n"
	for _, cat := range categories {
		if cat.Type == "Expense" {
			expense += fmt.Sprintf("%d. %s\n", cat.CategoryID, cat.Name)
		} else {
			income += fmt.Sprintf("%d. %s\n", cat.CategoryID, cat.Name)
		}
	}
	fmt.Print(expense)
	fmt.Print(income)
}

// findCategoryDetailByCid returns the name or type of
// a category given a category ID and requested detail.
func findCategoryDetailByCid(catID int, detail string) string {
	for _, cat := range categories {
		if catID == cat.CategoryID {
			if detail == "Type" {
				return cat.Type
			} else if detail == "Name" {
				return cat.Name
			}
		}
	}
	return ""
}

// generateReport handles sending a GET request to get all transactions
// then prints a summary report based on the retrieved user transactions.
func generateReport(c http.Client) {
	url := baseURL + "transactions"
	if trans, ok := getTransactions(c, url, 0); ok {
		if len(trans) > 0 {
			var totalIncome, totalExpense float64
			amountPerIncome := make(map[string]float64)
			amountPerExpense := make(map[string]float64)
			for _, tran := range trans {
				catType := findCategoryDetailByCid(tran.CategoryID, "Type")
				catName := findCategoryDetailByCid(tran.CategoryID, "Name")
				if catType == "Income" {
					totalIncome += tran.Amount
					amountPerIncome[catName] += tran.Amount
				} else {
					totalExpense += tran.Amount
					amountPerExpense[catName] += tran.Amount
				}
			}

			fmt.Println("========================")
			fmt.Println("     SUMMARY REPORT     ")
			fmt.Println("========================")
			fmt.Printf("Total Inflow: %s\n", utils.FormatFloat(totalIncome))
			fmt.Printf("Total Outflow: %s\n", utils.FormatFloat(totalExpense))
			fmt.Printf("Net Income: %s\n\n", utils.FormatFloat(totalIncome-totalExpense))
			fmt.Println("Income Transactions:")
			for k, v := range amountPerIncome {
				fmt.Printf("%s (%s%%) - %s\n", utils.FormatFloat(v), utils.FormatFloat((v/totalIncome)*100), k)
			}
			fmt.Println()
			fmt.Println("Expense Transactions:")
			for k, v := range amountPerExpense {
				fmt.Printf("%s (%s%%) - %s\n", utils.FormatFloat(v), utils.FormatFloat((v/totalExpense)*100), k)
			}
		} else {
			fmt.Println("No transaction/s found.")
		}
	}
}
