package main

import "fmt"

func main() {
	var n, counter int
	fmt.Scan(&n)
	for n != 0 {
		var stacks []int
		var h, reqH, moves int
		for i := 0; i < n; i++ {
			fmt.Scan(&h)
			reqH += h
			stacks = append(stacks, h)
		}
		reqH /= n
		for i := 0; i < n; i++ {
			var diff = stacks[i] - reqH
			if diff < 0 {
				moves += (-diff)
			}
		}
		counter++
		fmt.Printf("Set #%d\nThe minimum number of moves is %d.\n\n", counter, moves)
		fmt.Scan(&n)
	}
}
