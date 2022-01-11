package main

import "fmt"

func main() {
	var n, x, y, z int
	fmt.Scan(&n)
	for i := 1; i <= n; i++ {
		fmt.Scan(&x, &y, &z)
		if (x <= 20) && (y <= 20) && (z <= 20) {
			fmt.Printf("Case %d: %s\n", i, "good")
		} else {
			fmt.Printf("Case %d: %s\n", i, "bad")
		}
	}
}
