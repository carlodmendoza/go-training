package main

import (
	"fmt"
	"io"
)

func main() {
	var x, y int
	_, err := fmt.Scan(&x, &y)
	for {
		if err == io.EOF {
			break
		} else {
			if x > y {
				fmt.Println(x - y)
			} else {
				fmt.Println(y - x)
			}
			_, err = fmt.Scan(&x, &y)
		}
	}

}
