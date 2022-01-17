package main

import (
	"fmt"
	"io"
)

func main() {
	charMap := map[int]string{0: "B", 1: "G", 2: "C"}
	// possible permutations sorted alphabetically
	perms := [6][3]int{{0, 2, 1}, {0, 1, 2}, {2, 0, 1}, {2, 1, 0}, {1, 0, 2}, {1, 2, 0}}
	var num int
	_, err := fmt.Scan(&num)
	for {
		if err == io.EOF {
			break
		} else {
			var bins [3][3]int
			bins[0][0] = num
			// store inputs in 2d array
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if i == 0 && j == 0 {
						continue
					}
					_, err = fmt.Scan(&num)
					bins[i][j] = num
				}
			}

			var minPermIndices [3]int
			var minCount int
			// function that converts indices of perm to characters based on charMap
			indexToChar := func(perm [3]int) string {
				var charPerm string
				for _, char := range perm {
					charPerm += charMap[char]
				}
				return charPerm
			}
			// count # of movements needed per perm then return min perm and min count
			for i, perm := range perms {
				var count int
				for j, index := range perm {
					if j == 0 {
						count = count + bins[1][index] + bins[2][index]
					} else if j == 1 {
						count = count + bins[0][index] + bins[2][index]
					} else {
						count = count + bins[0][index] + bins[1][index]
					}
				}
				if i == 0 {
					minPermIndices = perm
					minCount = count
				} else {
					if count < minCount {
						minPermIndices = perm
						minCount = count
					}
				}
			}
			fmt.Println(indexToChar(minPermIndices), minCount)

			_, err = fmt.Scan(&num)
		}
	}
}
