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
		fmt.Printf("Case %d: %s", counter, words[word])
		fmt.Scan(&word)
		counter++
	}
}
