package etl

import "strings"

// Transform transforms data from map into an array data structure
// then returns the data in new format
func Transform(in map[int][]string) map[string]int {
	newMap := make(map[string]int)
	for k, v := range in {
		for _, letter := range v {
			newMap[strings.ToLower(letter)] = k
		}
	}
	return newMap
}
