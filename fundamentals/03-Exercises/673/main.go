package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	pairs := map[string]string{")": "(", "]": "["}
	// create scanner to check for empty lines
	scanner := bufio.NewScanner(os.Stdin)
	// skip first line
	scanner.Scan()
	for scanner.Scan() {
		var stack []string
		line := scanner.Text()
		for _, char := range line {
			// if opening character, push to the stack; else, pop last element of stack if stack is not empty and popped element pairs with the closing character
			if string(char) == "(" || string(char) == "[" {
				stack = append(stack, string(char))
			} else {
				// if closing character, and if stack is empty or popped element does not pair with closing character, just push to the stack
				if len(stack) > 0 {
					lastElem := string(stack[len(stack)-1])
					if lastElem == pairs[string(char)] {
						stack = stack[:len(stack)-1]
					} else {
						stack = append(stack, string(char))
						break
					}
				} else {
					stack = append(stack, string(char))
					break
				}
			}
		}
		// at the end, if stack is empty, parentheses are balanced; else, they are not balanced
		if len(stack) == 0 {
			fmt.Println("Yes")
		} else {
			fmt.Println("No")
		}
	}
}
