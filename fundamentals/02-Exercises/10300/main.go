package main

import "fmt"

func main() {
	var n, farmers, size, animals, ecoScore, premium, sum int
	fmt.Scan(&n)
	for i := 1; i <= n; i++ {
		fmt.Scan(&farmers)
		sum = 0
		for j := 1; j <= farmers; j++ {
			fmt.Scan(&size, &animals, &ecoScore)
			premium = size * ecoScore
			sum += premium
		}
		fmt.Println(sum)
	}
}
