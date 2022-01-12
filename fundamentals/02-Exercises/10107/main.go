package main

import (
	"fmt"
	"io"
	"sort"
)

func main() {
	var x []int
	var num int
	_, err := fmt.Scan(&num)
	for {
		if err == io.EOF {
			break
		} else {
			x = append(x, num)
			fmt.Println(getMedian(x))
			_, err = fmt.Scan(&num)
		}
	}
}

func getMedian(arr []int) int {
	sort.Ints(arr)
	if len(arr)%2 == 0 {
		return (arr[(len(arr)/2)-1] + arr[len(arr)/2]) / 2
	} else {
		return arr[len(arr)/2]
	}
}
