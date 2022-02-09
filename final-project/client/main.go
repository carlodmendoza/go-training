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

const baseURL = "http://localhost:8080/"

var cookie *http.Cookie
var categories []models.Category

func main() {
	var choice int
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	utils.PrintWelcomeMessage()
	commands := []string{"Sign in", "Sign up"}
	utils.PrintValidCommands(commands)
	fmt.Scan(&choice)
	for {
		if choice == 1 {
			if signin(c) {
				break
			} else {
				utils.PrintValidCommands(commands)
				fmt.Scan(&choice)
			}
		} else if choice == 2 {
			signup(c)
			utils.PrintValidCommands(commands)
			fmt.Scan(&choice)
		} else {
			utils.PrintValidCommands(commands)
			fmt.Scan(&choice)
		}
	}

	commands = []string{"View my transactions", "View report", "Add new transaction", "View a transaction", "Edit a transaction", "Delete a transaction", "Delete all transactions"}
	categories = getCategories(c)
	utils.PrintValidCommands(commands)
	fmt.Scan(&choice)
	for {
		if choice == 1 {
			var transactions []models.Transaction
			viewTransactions(c, transactions, 0)
			utils.PrintValidCommands(commands)
			fmt.Scan(&choice)
		} else if choice == 2 {
			utils.PrintValidCommands(commands)
			fmt.Scan(&choice)
		} else if choice == 3 {
			utils.PrintValidCommands(commands)
			fmt.Scan(&choice)
		} else if choice == 4 {
			var transaction models.Transaction
			var id int
			fmt.Println("Enter the transaction ID: ")
			fmt.Scan(&id)
			viewTransactions(c, transaction, id)
			utils.PrintValidCommands(commands)
			fmt.Scan(&choice)
		} else if choice == 5 {
			utils.PrintValidCommands(commands)
			fmt.Scan(&choice)
		} else if choice == 6 {
			utils.PrintValidCommands(commands)
			fmt.Scan(&choice)
		} else if choice == 7 {
			utils.PrintValidCommands(commands)
			fmt.Scan(&choice)
		} else {
			utils.PrintValidCommands(commands)
			fmt.Scan(&choice)
		}
	}
}

func signin(c http.Client) bool {
	url := baseURL + "signin"

	var username, password string
	fmt.Println("Enter your username: ")
	fmt.Scan(&username)
	fmt.Println("Enter your password: ")
	fmt.Scan(&password)

	reqBody := fmt.Sprintf("{\"username\":\"%s\", \"password\":\"%s\"}", username, password)
	return getResponse(c, url, "POST", reqBody)
}

func signup(c http.Client) bool {
	url := baseURL + "signup"

	var name, username, password string
	fmt.Println("Enter your name: ")
	fmt.Scan(&name)
	fmt.Println("Enter your username: ")
	fmt.Scan(&username)
	fmt.Println("Enter your password: ")
	fmt.Scan(&password)

	reqBody := fmt.Sprintf("{\"name\":\"%s\", \"username\":\"%s\", \"password\":\"%s\"}", name, username, password)
	return getResponse(c, url, "POST", reqBody)
}

func viewTransactions(c http.Client, model interface{}, transID int) {
	var url string
	switch model.(type) {
	case []models.Transaction:
		url = baseURL + "transactions"
	case models.Transaction:
		url = baseURL + "transactions/" + fmt.Sprint(transID)
	}
	if trans, ok := getTransactions(c, url, model); ok {
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

func getTransactions(c http.Client, url string, model interface{}) ([]models.Transaction, bool) {
	req, _ := http.NewRequest("GET", url, nil)
	req.AddCookie(cookie)
	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("Failed to get response: %s", err.Error())
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	switch v := model.(type) {
	case []models.Transaction:
		if err := json.Unmarshal(data, &v); err != nil {
			fmt.Printf("Failed to parse json response: %s\n", err.Error())
			return []models.Transaction{}, false
		}
		return v, true
	case models.Transaction:
		if err := json.Unmarshal(data, &v); err != nil {
			fmt.Printf("Failed to parse json response: %s\n", err.Error())
			return []models.Transaction{}, false
		}
		// if no transaction found
		if v.TransactionID == 0 {
			return []models.Transaction{}, true
		}
		return []models.Transaction{v}, true
	default:
		return []models.Transaction{}, false
	}
}

func getResponse(c http.Client, url, method, reqBody string) bool {
	var response models.Response
	var resp *http.Response
	var err error

	switch method {
	case "POST":
		outData := bytes.NewBuffer([]byte(reqBody))
		resp, err = c.Post(url, "application/json", outData)
	}

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

func getCategories(c http.Client) []models.Category {
	url := baseURL + "categories"
	var categories []models.Category
	var resp *http.Response
	var err error

	req, _ := http.NewRequest("GET", url, nil)
	req.AddCookie(cookie)
	resp, err = c.Do(req)
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

func getCategoryDetails(catID int) string {
	for _, cat := range categories {
		if catID == cat.CategoryID {
			return fmt.Sprintf("%s (%s)", cat.Name, cat.Type)
		}
	}
	return fmt.Sprint(catID)
}
