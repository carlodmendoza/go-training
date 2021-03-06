package main

import "fmt"

func main() {
	var n int
	var x, y, z int32
	fmt.Scan(&n)
	for i := 1; i <= n; i++ {
		fmt.Scan(&x, &y, &z)
		if x+y > z && x+z > y && y+z > x {
			if x == y && y == z {
				fmt.Printf("Case %d: Equilateral\n", i)
			} else if x == y || y == z || x == z {
				fmt.Printf("Case %d: Isosceles\n", i)
			} else {
				fmt.Printf("Case %d: Scalene\n", i)
			}
		} else {
			fmt.Printf("Case %d: Invalid\n", i)
		}
	}
}
