package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	var x, y, z int
	var arr = []int{}
	fmt.Scan(&x, &y, &z)
	for x != 0 && y != 0 && z != 0 {
		arr = []int{x, y, z}
		sort.Ints(arr)
		if math.Pow(float64(arr[0]), 2)+math.Pow(float64(arr[1]), 2) == math.Pow(float64(arr[2]), 2) {
			fmt.Println("right")
		} else {
			fmt.Println("wrong")
		}
		fmt.Scan(&x, &y, &z)
	}
}
