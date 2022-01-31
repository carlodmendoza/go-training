package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
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
	for _, question := range questions {
		if question[0] == "" || question[1] == "" {
			log.Fatal("Incorrect database format")
		}
	}

	// print every question then get answer from input
	// output number of correct answers at the end
	var count int
	var finishedQuestions [][]string
	for i := 1; i <= *n; i++ {
		// if all questions are already answered, get the questions again
		if len(questions) == 0 {
			questions = finishedQuestions
			finishedQuestions = nil
		}
		// generate random index from 0 to len(questions)-1
		rand.Seed(time.Now().UnixNano())
		questionIndex := rand.Intn(len(questions))
		// print question then get input. If answer is correct, increment counter
		var ans string
		fmt.Printf("Q%d: %s = ", i, questions[questionIndex][0])
		fmt.Scan(&ans)
		// remove leading and trailing whitespaces
		questions[questionIndex][1] = strings.TrimSpace(questions[questionIndex][1])
		if questions[questionIndex][1] == ans {
			count++
		}
		// add question to a different array for retrieval later
		finishedQuestions = append(finishedQuestions, questions[questionIndex])
		// remove question from original array
		questions = removeElementFromArray(questions, questionIndex)
	}
	fmt.Printf("You answered %d out of %d questions correctly.", count, *n)
}

// removeElementFromArray removes an element given an index
// by replacing it with the last element and returning a copy
// that excludes the last element.
func removeElementFromArray(s [][]string, i int) [][]string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
