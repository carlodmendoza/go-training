package main

import "fmt"

func main() {
	var scores []float32
	var n, students, counter int
	var score, aveScore, aboveAve float32
	fmt.Scan(&n)
	for i := 1; i <= n; i++ {
		fmt.Scan(&students)
		scores = []float32{}
		aveScore = 0.0
		counter = 0
		aboveAve = 0.0
		for j := 1; j <= students; j++ {
			fmt.Scan(&score)
			scores = append(scores, score)
			aveScore += score
		}
		aveScore = aveScore / float32(students)
		for _, v := range scores {
			if v > aveScore {
				counter++
			}
		}
		aboveAve = (float32(counter) / float32(students)) * 100
		fmt.Printf("%.3f%%\n", aboveAve)
	}
}
