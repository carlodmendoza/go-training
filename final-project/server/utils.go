package main

import (
	"fmt"
	"net/http"
)

func sendMessage(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	responseBody := fmt.Sprintf("{\"message\": \"%s\"}", message)
	fmt.Fprintln(w, responseBody)
}

func sendMessageWithBody(w http.ResponseWriter, success bool, message string) {
	w.Header().Set("Content-Type", "application/json")
	responseBody := fmt.Sprintf(
		"{\"success\": %t,"+
			"\"message\": \"%s\"}",
		success, message)
	fmt.Fprintln(w, responseBody)
}
