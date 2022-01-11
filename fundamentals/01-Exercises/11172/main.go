package main

import (
	"fmt"
)

func main() {
	var n, x, y int
	fmt.Scan(&n)
	for i := 1; i <= n; i++ {
		fmt.Scan(&x, &y)
		if x > y {
			fmt.Println(">")
		} else if x < y {
			fmt.Println("<")
		} else {
			fmt.Println("=")
		}
	}
}
