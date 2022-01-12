package main

import "fmt"

func main() {
	var n int
	var word string
	fmt.Scan(&n)
	for i := 1; i <= n; i++ {
		fmt.Scan(&word)
		if len(word) == 5 {
			fmt.Println(3)
		} else {
			if (string(word[0]) == "o" && string(word[1]) == "n") || (string(word[0]) == "o" && string(word[2]) == "e") || (string(word[1]) == "n" && string(word[2]) == "e") {
				fmt.Println(1)
			} else {
				fmt.Println(2)
			}
		}
	}
}
