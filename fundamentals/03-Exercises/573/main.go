package main

import "fmt"

func main() {
	var H, U, D, F float32
	fmt.Scan(&H)
	for H != 0 {
		fmt.Scan(&U, &D, &F)
		var counter int
		var currHeight float32
		lostDistance := (F / 100) * U
		for {
			counter++
			if U > 0 {
				currHeight += U
			}
			if currHeight > float32(H) {
				fmt.Printf("success on day %d\n", counter)
				break
			}
			currHeight -= D
			if currHeight < 0 {
				fmt.Printf("failure on day %d\n", counter)
				break
			}
			U -= lostDistance
		}
		fmt.Scan(&H)
	}
}
