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
	"time"
)

const baseURL = "http://localhost:8080/"

var cookie *http.Cookie

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
	utils.PrintValidCommands(commands)
	fmt.Scan(&choice)
	for {
		if choice == 1 {
			viewAllTransactions(c)
			utils.PrintValidCommands(commands)
			fmt.Scan(&choice)
		} else if choice == 2 {
			utils.PrintValidCommands(commands)
			fmt.Scan(&choice)
		} else if choice == 3 {
			utils.PrintValidCommands(commands)
			fmt.Scan(&choice)
		} else if choice == 4 {
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

func viewAllTransactions(c http.Client) {
	url := baseURL + "transactions"
	if trans, ok := getTransactionFromResponse(c, url); ok {
		if len(trans) > 0 {
			fmt.Println(trans)
		} else {
			fmt.Println("No transactions found.")
		}
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

func getTransactionFromResponse(c http.Client, url string) ([]models.Transaction, bool) {
	var response models.Response
	var transactions []models.Transaction
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
	if resp.StatusCode == http.StatusUnauthorized {
		if err := json.Unmarshal(data, &response); err != nil {
			fmt.Printf("Failed to parse json response: %s\n", err.Error())
			return []models.Transaction{}, false
		}
		fmt.Printf("Error: %s\n", response.Message)
		return []models.Transaction{}, false
	}
	if err := json.Unmarshal(data, &transactions); err != nil {
		fmt.Printf("Failed to parse json response: %s\n", err.Error())
		return []models.Transaction{}, false
	}
	return transactions, true
}
