package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
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

// func (wc *WordCounter) Value(key string) int {
// 	wc.mu.Lock()
// 	defer wc.mu.Unlock()
// 	return wc.counts[key]
// }

func main() {
	counts := make(map[string]int)
	ch := make(chan map[string]int)
	var keys []string
	// process the files concurrently
	for _, fileName := range os.Args[1:] {
		go processFile(fileName, ch)
		// receive the word counts from channel then update main counter
		for k, v := range <-ch {
			counts[k] += v
		}
	}
	// list the keys in alphabetical order then print it and its corresponding value
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k, counts[k])
	}
}

func processFile(fileName string, ch chan map[string]int) {
	wc := WordCounter{counts: make(map[string]int)}
	// open file
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// read file
	scanner := bufio.NewScanner(file)
	// get all words separated by spaces
	scanner.Split(bufio.ScanWords)
	// for each word, get only valid characters (letters or digits)
	for scanner.Scan() {
		var word string
		for _, char := range scanner.Text() {
			if unicode.IsDigit(char) || unicode.IsLetter(char) {
				word += string(char)
			}
		}
		// increment count of word if it's not an empty string
		if word != "" {
			wc.Inc(word)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to parse file: %s", err)
	}
	// send the word counts to the channel
	ch <- wc.counts
}
