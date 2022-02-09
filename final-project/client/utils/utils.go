package utils

import (
	"fmt"
	"strconv"
)

func PrintWelcomeMessage() {
	fmt.Println("========================================")
	fmt.Println("Welcome to your Personal Budget Tracker!")
	fmt.Println("========================================")
}

func PrintExitMessage() {
	fmt.Println("===========================================")
	fmt.Println("Thank you for using the app. See you again!")
	fmt.Println("===========================================")
}

func PrintCommands(commands []string) {
	fmt.Println("\nWhat do you want to do? (Enter the number of your choice):")
	counter := 1
	for _, cmd := range commands {
		fmt.Printf("%d. %s\n", counter, cmd)
		counter++
	}
}

func FormatFloat(num float64) string {
	if num == float64(int(num)) {
		return strconv.Itoa(int(num))
	}
	return strconv.FormatFloat(num, 'f', 2, 64)
}
