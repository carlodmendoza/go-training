package main

import (
	"fmt"
	"io"
)

func main() {
	var v, t int
	_, err := fmt.Scan(&v, &t)
	for {
		if err == io.EOF {
			break
		} else {
			fmt.Println(v * 2 * t)
			_, err = fmt.Scan(&v, &t)
		}
	}
}
