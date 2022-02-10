package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

/*
	Utils program contains utility functions for sending
	responses to a client and generating a token.
	Author: Carlo Mendoza
*/

// SendMessage writes a formatted JSON to a client given
// a message field.
func SendMessage(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	responseBody := fmt.Sprintf("{\"message\": \"%s\"}", message)
	fmt.Fprintln(w, responseBody)
}

// SendMessageWithBody writes a formatted JSON to a client given
// a success and message field.
func SendMessageWithBody(w http.ResponseWriter, success bool, message string) {
	w.Header().Set("Content-Type", "application/json")
	responseBody := fmt.Sprintf(
		"{\"success\": %t,"+
			"\"message\": \"%s\"}",
		success, message)
	fmt.Fprintln(w, responseBody)
}

// GenerateSessionToken returns a token for authorizing
// client requests.
func GenerateSessionToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
