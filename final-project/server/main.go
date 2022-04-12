package main

import (
	"fmt"
	"log"
	"net/http"
	"server/auth"
	"server/categories"
	"server/storage"
	"server/storage/filebased"
	"server/transactions"
)

func main() {
	fmt.Println("Server running on port 8080")
	err := http.ListenAndServe(":8080", handler(filebased.FileDB))
	if err != nil {
		log.Fatalf("Error ListenAndServe(): %s", err.Error())
	}
}

// handler handles requests to the server depending on the request URL given a StorageService.
// It authorizes a user to make requests given a request token.
func handler(db storage.StorageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var transID int
		if r.URL.Path == "/signin" {
			auth.Signin(db, w, r)
		} else if r.URL.Path == "/signup" {
			auth.Signup(db, w, r)
		} else if r.URL.Path == "/transactions" {
			uid := auth.AuthenticateToken(db, w, r)
			if uid <= 0 {
				return
			}
			transactions.ProcessTransaction(db, w, r, uid)
		} else if n, _ := fmt.Sscanf(r.URL.Path, "/transactions/%d", &transID); n == 1 {
			uid := auth.AuthenticateToken(db, w, r)
			if uid <= 0 {
				return
			}
			transactions.ProcessTransactionID(db, w, r, uid, transID)
		} else if r.URL.Path == "/categories" {
			uid := auth.AuthenticateToken(db, w, r)
			if uid <= 0 {
				return
			}
			categories.ProcessCategories(db, w, r)
		} else {
			http.Error(w, "Invalid URL or request", http.StatusNotImplemented)
		}
	}
}
