package utils

import (
	"fmt"
	"net/http"
)

func SendMessage(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	responseBody := fmt.Sprintf("{\"message\": \"%s\"}", message)
	fmt.Fprintln(w, responseBody)
}

func SendMessageWithBody(w http.ResponseWriter, success bool, message string) {
	w.Header().Set("Content-Type", "application/json")
	responseBody := fmt.Sprintf(
		"{\"success\": %t,"+
			"\"message\": \"%s\"}",
		success, message)
	fmt.Fprintln(w, responseBody)
}
