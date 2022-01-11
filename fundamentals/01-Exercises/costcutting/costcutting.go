package main

import (
	"fmt"
	"sort"
)

func main() {
	var n, x, y, z int
	var arr = []int{}
	fmt.Scan(&n)
	for i := 1; i <= n; i++ {
		fmt.Scan(&x, &y, &z)
		arr = []int{x, y, z}
		sort.Ints(arr)
		fmt.Printf("Case %d: %d", i, arr[1])
	}
}
