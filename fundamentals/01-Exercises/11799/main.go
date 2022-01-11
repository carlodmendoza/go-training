package main

import (
	"fmt"
)

func main() {
	var n int
	fmt.Scan(&n)
	for i := 1; i <= n; i++ {
		var count int
		fmt.Scan(&count)
		var max, next int
		fmt.Scan(&max)
		for j := 1; j <= count-1; j++ {
			fmt.Scan(&next)
			if next > max {
				max = next
			}
		}
		fmt.Printf("Case %d: %d\n", i, max)
	}
}
