package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/storage/filebased"
	"syscall"
)

func main() {
	fmt.Println("Server running on port 8080")

	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) //nolint

	// TODO: create env var for file path and chosen storage service
	storage := filebased.Initialize("../deploy/dev/server/storage/data")
	r := GetRouter(storage)

	go func() {
		<-sigChannel
		err := storage.Shutdown()
		if err != nil {
			fmt.Printf("Shutdown error: %s\n", err)
		}
		log.Fatalf("Shutting down the server")
	}()

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Error ListenAndServe(): %s", err)
	}
}

// handler handles requests to the server depending on the request URL given a StorageService.
// It authorizes a user to make requests given a request token.
// func handler(db storage.Service) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var transID int
// 		if r.URL.Path == "/signin" {
// 			auth.Signin(db, w, r)
// 		} else if r.URL.Path == "/signup" {
// 			auth.Signup(db, w, r)
// 		} else if r.URL.Path == "/transactions" {
// 			user := auth.AuthenticateToken(db, w, r)
// 			if user == "" {
// 				return
// 			}
// 			transactions.ProcessTransaction(db, w, r, user)
// 		} else if n, _ := fmt.Sscanf(r.URL.Path, "/transactions/%d", &transID); n == 1 {
// 			user := auth.AuthenticateToken(db, w, r)
// 			if user == "" {
// 				return
// 			}
// 			transactions.ProcessTransactionID(db, w, r, user, transID)
// 		} else if r.URL.Path == "/categories" {
// 			user := auth.AuthenticateToken(db, w, r)
// 			if user == "" {
// 				return
// 			}
// 			categories.ProcessCategories(db, w, r)
// 		} else {
// 			http.Error(w, "Invalid URL or request", http.StatusNotImplemented)
// 		}
// 	}
// }
