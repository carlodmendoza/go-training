package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Message string
	Success bool
}

const baseURL = "http://localhost:8080/"

var token string

func main() {
	printWelcomeMessage()
	var choice int
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	commands := []string{"Sign in", "Sign up"}
	printValidCommands(commands)
	fmt.Scan(&choice)
	for {
		if choice == 1 {
			if signin(c) {
				break
			} else {
				printValidCommands(commands)
				fmt.Scan(&choice)
			}
		} else if choice == 2 {
			signup(c)
			printValidCommands(commands)
			fmt.Scan(&choice)
		} else {
			printValidCommands(commands)
			fmt.Scan(&choice)
		}
	}
}

func printWelcomeMessage() {
	fmt.Println("========================================")
	fmt.Println("Welcome to your Personal Budget Tracker!")
	fmt.Println("========================================")
}

func printValidCommands(commands []string) {
	fmt.Println("\nWhat do you want to do? (Enter the number of your choice):")
	counter := 1
	for _, cmd := range commands {
		fmt.Printf("%d. %s\n", counter, cmd)
		counter++
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

func getResponse(c http.Client, url, method, reqBody string) bool {
	var response Response
	var resp *http.Response
	var err error
	outData := bytes.NewBuffer([]byte(reqBody))

	switch method {
	case "POST":
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
	} else {
		if response.Success {
			if url == baseURL+"signin" {
				token = resp.Cookies()[0].Value
			}
			fmt.Println(response.Message)
		} else {
			fmt.Printf("Error: %s\n", response.Message)
		}
		return response.Success
	}
}
