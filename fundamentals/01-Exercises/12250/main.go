package main

import "fmt"

func main() {
	words := map[string]string{
		"HELLO":        "ENGLISH",
		"HOLA":         "SPANISH",
		"HALLO":        "GERMAN",
		"BONJOUR":      "FRENCH",
		"CIAO":         "ITALIAN",
		"ZDRAVSTVUJTE": "RUSSIAN",
	}
	var word string
	var counter int = 1
	fmt.Scan(&word)
	for word != "#" {
		if val, exists := words[word]; exists {
			fmt.Printf("Case %d: %s\n", counter, val)
		} else {
			fmt.Printf("Case %d: UNKNOWN\n", counter)
		}
		fmt.Scan(&word)
		counter++
	}
}
