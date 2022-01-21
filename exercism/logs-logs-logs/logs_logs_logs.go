package logs

import (
    "fmt"
    "unicode/utf8"
)

// Application identifies the application emitting the given log.
func Application(log string) string {
	var sourceApp string
	for _, char := range log {
		unicode := fmt.Sprintf("%U", char)
		if unicode == "U+2757" {
			sourceApp = "recommendation"
			break
		} else if unicode == "U+1F50D" {
			sourceApp = "search"
			break
		} else if unicode == "U+2600" {
			sourceApp = "weather"
			break
		} else {
			sourceApp = "default"
		}
	}
	return sourceApp
}

// Replace replaces all occurrences of old with new, returning the modified log
// to the caller.
func Replace(log string, oldRune, newRune rune) string {
	var newString string
    oldRuneUnicode := fmt.Sprintf("%U", oldRune)
    for _, char := range log {
        unicode := fmt.Sprintf("%U", char)
        if unicode == oldRuneUnicode {
            newString += string(newRune)
        } else {
            newString += string(char)
        }
    }
	return newString
}

// WithinLimit determines whether or not the number of characters in log is
// within the limit.
func WithinLimit(log string, limit int) bool {
    return utf8.RuneCountInString(log) <= limit
}
