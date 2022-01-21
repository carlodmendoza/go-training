package etl

import "strings"

func Transform(in map[int][]string) map[string]int {
	newMap := make(map[string]int)
	for k, v := range in {
		for _, letter := range v {
			newMap[strings.ToLower(letter)] = k
		}
	}
	return newMap
}
