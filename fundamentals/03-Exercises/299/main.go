package main

import "fmt"

func main() {
	var N, L int
	fmt.Scan(&N)
	for i := 1; i <= N; i++ {
		fmt.Scan(&L)
		var num, swapCount int
		var order []int
		// store inputs in array
		for j := 0; j < L; j++ {
			fmt.Scan(&num)
			order = append(order, num)
		}
		// bubble sort array and return # of swaps done
		for index := 0; index < len(order); index++ {
			for elemIndex := 0; elemIndex < L-index-1; elemIndex++ {
				if order[elemIndex] > order[elemIndex+1] {
					temp := order[elemIndex]
					order[elemIndex] = order[elemIndex+1]
					order[elemIndex+1] = temp
					swapCount++
				}
			}
		}
		fmt.Printf("Optimal train swapping takes %d swaps.\n", swapCount)
	}
}
