package utils

import (
	"fmt"
)

func PrintWelcomeMessage() {
	fmt.Println("========================================")
	fmt.Println("Welcome to your Personal Budget Tracker!")
	fmt.Println("========================================")
}

func PrintValidCommands(commands []string) {
	fmt.Println("\nWhat do you want to do? (Enter the number of your choice):")
	counter := 1
	for _, cmd := range commands {
		fmt.Printf("%d. %s\n", counter, cmd)
		counter++
	}
}
