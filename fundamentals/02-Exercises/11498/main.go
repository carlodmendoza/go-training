package main

import "fmt"

func main() {
	var n, originX, originY, x, y int
	fmt.Scan(&n)
	for n != 0 {
		fmt.Scan(&originX, &originY)
		for i := 1; i <= n; i++ {
			fmt.Scan(&x, &y)
			if x == originX || y == originY {
				fmt.Println("divisa")
			} else if x < originX && y > originY {
				fmt.Println("NO")
			} else if x > originX && y > originY {
				fmt.Println("NE")
			} else if x > originX && y < originY {
				fmt.Println("SE")
			} else {
				fmt.Println("SO")
			}
		}
		fmt.Scan(&n)
	}

}
