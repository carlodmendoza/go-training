package main

import "fmt"

func main() {
	words := map[string]string{
		"Hajj":  "Hajj-e-Akbar",
		"Umrah": "Hajj-e-Asghar",
	}
	var word string
	var counter int = 1
	fmt.Scan(&word)
	for word != "*" {
		fmt.Printf("Case %d: %s\n", counter, words[word])
		fmt.Scan(&word)
		counter++
	}
}
