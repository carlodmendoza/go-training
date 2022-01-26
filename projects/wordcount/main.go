package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
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

func (wc *WordCounter) Value(key string) int {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	return wc.counts[key]
}

func main() {
	wc := WordCounter{counts: make(map[string]int)}
	ch := make(chan []string)
	var wordList []string
	var keys []string
	// process the files concurrently
	for _, fileName := range os.Args[1:] {
		go processFile(fileName, ch)
	}
	// receive the list of words from channel then append to main list
	for range os.Args[1:] {
		wordList = append(wordList, <-ch...)
	}
	// count the words concurrently
	for _, word := range wordList {
		go wc.Inc(word)
	}
	// list the keys in alphabetical order then print it and its corresponding value
	for k := range wc.counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k, wc.Value(k))
	}
}

func processFile(fileName string, ch chan []string) {
	var words []string
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
	// for each word, get only valid characters (letters or digits) then convert to lower case
	for scanner.Scan() {
		var word string
		for _, char := range scanner.Text() {
			if unicode.IsDigit(char) || unicode.IsLetter(char) {
				word += string(char)
			}
		}
		word = strings.ToLower(word)
		// append word to array if it's not an empty string
		if word != "" {
			words = append(words, word)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to parse file: %s", err)
	}
	// send the word array to the channel
	ch <- words
}
