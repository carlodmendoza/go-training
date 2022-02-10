package utils

import (
	"fmt"
	"strconv"
)

/*
	Utils program contains utility functions for printing messages
	and commands, and formatting numbers.
	Author: Carlo Mendoza
*/

// PrintWelcomeMessage prints a welcome message upon app startup.
func PrintWelcomeMessage() {
	fmt.Println("========================================")
	fmt.Println("Welcome to your Personal Budget Tracker!")
	fmt.Println("========================================")
}

// PrintExitMessage prints an exit message upon app exit.
func PrintExitMessage() {
	fmt.Println("===========================================")
	fmt.Println("Thank you for using the app. See you again!")
	fmt.Println("===========================================")
}

// PrintCommands prints all valid commands that user can make
// given a list of string.
func PrintCommands(commands []string) {
	fmt.Println("\nWhat do you want to do? (Enter the number of your choice):")
	counter := 1
	for _, cmd := range commands {
		fmt.Printf("%d. %s\n", counter, cmd)
		counter++
	}
}

// FormatFloat returns the formatted string of a
// given floating-point number.
func FormatFloat(num float64) string {
	if num == float64(int(num)) {
		return strconv.Itoa(int(num))
	}
	return strconv.FormatFloat(num, 'f', 2, 64)
}
