package main

import "fmt"

func main() {
	var x, y int
	fmt.Scan(&x, &y)
	for x != -1 && y != -1 {
		if x < y {
			if (x - y + 100) < (y - x) {
				fmt.Println(x - y + 100)
			} else {
				fmt.Println(y - x)
			}
		} else if x > y {
			if (x - y) < (y - x + 100) {
				fmt.Println(x - y)
			} else {
				fmt.Println(y - x + 100)
			}
		} else {
			fmt.Println(0)
		}
		fmt.Scan(&x, &y)
	}
}
