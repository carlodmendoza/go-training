package main

import "fmt"

func main() {
	var cases, n, m int
	fmt.Scan(&cases)
	for i := 1; i <= cases; i++ {
		fmt.Scan(&n, &m)
		fmt.Println((n / 3) * (m / 3))
	}
}
