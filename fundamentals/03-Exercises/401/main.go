package main

import (
	"fmt"
	"io"
)

func main() {
	var word string
	reverses := map[string]string{"A": "A", "E": "3", "H": "H", "I": "I", "J": "L",
		"L": "J", "M": "M", "O": "O", "S": "2", "T": "T",
		"U": "U", "V": "V", "W": "W", "X": "X", "Y": "Y",
		"Z": "5", "1": "1", "2": "S", "3": "E", "5": "Z",
		"8": "8"}
	_, err := fmt.Scan(&word)
	for {
		var pal, rev string
		if err == io.EOF {
			break
		} else {
			for i := len(word) - 1; i >= 0; i-- {
				pal += string(word[i])
				rev += reverses[string(word[i])]
			}
			if pal == word && rev == word {
				fmt.Printf("%s -- is a mirrored palindrome.\n\n", word)
			} else if pal == word && rev != word {
				fmt.Printf("%s -- is a regular palindrome.\n\n", word)
			} else if pal != word && rev == word {
				fmt.Printf("%s -- is a mirrored string.\n\n", word)
			} else {
				fmt.Printf("%s -- is not a palindrome.\n\n", word)
			}
		}
		_, err = fmt.Scan(&word)
	}
}
