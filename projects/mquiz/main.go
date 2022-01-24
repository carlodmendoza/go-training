package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	// define flags for file name and number of questions
	fileName := flag.String("csv", "problem.csv", "file name of database")
	n := flag.Int("n", 10, "number of questions to get from database")
	flag.Parse()

	// open database file
	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatalf("Failed to open csv file: %s", err)
	}
	defer file.Close()

	// read database file
	reader := csv.NewReader(file)
	questions, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to parse csv file: %s", err)
	}

	// report errors
	if *n > len(questions) {
		fmt.Println("warning: insufficient questions, some questions will be repeated")
	}

	// print every question then get answer from input
	// output number of correct answers at the end
	var count int
	var finishedQuestions [][]string
	for i := 1; i <= *n; i++ {
		if len(questions) == 0 {
			questions = finishedQuestions
			finishedQuestions = nil
		}
		rand.Seed(time.Now().UnixNano())
		questionIndex := rand.Intn(len(questions))

		var ans string
		fmt.Printf("Q%d: %s = ", i, questions[questionIndex][0])
		fmt.Scan(&ans)
		if questions[questionIndex][1] == ans {
			count++
		}

		finishedQuestions = append(finishedQuestions, questions[questionIndex])
		questions = removeElementFromArray(questions, questionIndex)
	}
	fmt.Printf("You answered %d out of %d questions correctly.", count, *n)
}

func removeElementFromArray(s [][]string, i int) [][]string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
