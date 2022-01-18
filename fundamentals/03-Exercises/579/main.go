package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	var time string
	fmt.Scan(&time)
	for time != "0:00" {
		var hourDeg, minDeg, angle float64
		hour, _ := strconv.ParseFloat(strings.Split(time, ":")[0], 64)
		min, _ := strconv.ParseFloat(strings.Split(time, ":")[1], 64)
		hourDeg = hour*30 + ((min / 60) * 30)
		minDeg = min * 6
		angle = math.Abs(hourDeg - minDeg)
		if angle > 180 {
			angle = 360 - angle
		}
		fmt.Printf("%.3f\n", angle)
		fmt.Scan(&time)
	}
}
