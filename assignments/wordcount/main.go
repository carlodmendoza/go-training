package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
)

type WordCounter struct {
	mu     sync.Mutex
	counts map[string]int
}

func (wc *WordCounter) Inc(key string) {
	wc.mu.Lock()
	wc.counts[key]++
	wc.mu.Unlock()
}

func (wc *WordCounter) Value(key string) int {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	return wc.counts[key]
}

func main() {
	wc := WordCounter{counts: make(map[string]int)}
	fileWordCountChan := make(chan map[string]int)
	var keys []string

	// process the files concurrently
	for _, fileName := range os.Args[1:] {
		go countWordsInFile(fileName, fileWordCountChan)
	}
	// receive the word count of each file from the channel then count the words concurrently
	for range os.Args[1:] {
		fileWordCount := <-fileWordCountChan
		for k := range fileWordCount {
			go wc.Inc(k)
		}
	}
	// delay before main routine prints the output
	time.Sleep(time.Second)

	// list the keys in alphabetical order then print it and its corresponding value
	for k := range wc.counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k, wc.Value(k))
	}
}

// countWordsInFile processes a file by getting all valid words from the file
// then sending a map of the word counts to a map channel
func countWordsInFile(fileName string, ch chan map[string]int) {
	var counter int
	wordCount := make(map[string]int)
	wordChan := make(chan string)

	// open file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// read file then get all words separated by spaces
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	// process each word concurrently
	for scanner.Scan() {
		go processWord(scanner.Text(), wordChan)
		counter++
	}
	// receive the new word from the channel then count it if it's not an empty string
	for i := 0; i < counter; i++ {
		word := <-wordChan
		if word != "" {
			wordCount[word] += 1
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to parse file: %s", err)
	}

	ch <- wordCount
}

// processWord removes non-letter and non-digit characters from a string
// then sends the new string in lowercase to a string channel
func processWord(oldWord string, ch chan string) {
	var newWord string
	for _, char := range oldWord {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			newWord += string(char)
		}
	}
	newWord = strings.ToLower(newWord)
	ch <- newWord
}
