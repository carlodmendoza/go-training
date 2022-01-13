package main

import "fmt"

func main() {
	var n, a, b, sum int
	fmt.Scan(&n)
	for i := 1; i <= n; i++ {
		fmt.Scan(&a, &b)
		sum = 0
		for j := a; j <= b; j++ {
			if j%2 == 1 {
				sum += j
			}
		}
		fmt.Printf("Case %d: %d\n", i, sum)
	}
}
