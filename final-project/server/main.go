package main

import (
	"final-project/server/models"
	"final-project/server/utils"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Server running in port 8080")
	db := &models.Database{}
	db.InitializeDatabase()
	if err := http.ListenAndServe("localhost:8080", handler(db)); err != nil {
		log.Fatalf("Error ListenAndServe(): %s", err.Error())
	}
}

func handler(db *models.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/signin" {
			db.Signin(w, r)
		} else if r.URL.Path == "/signup" {
			db.Signup(w, r)
		} else if r.URL.Path == "/transactions" {
			if db.CurrentUser.UserID == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				utils.SendMessageWithBody(w, false, "Please login first.")
			} else {
				// VIEW transactions or transaction
				// POST transaction
				// PUT transaction
				// DELETE transaction
				fmt.Println("Handle transactions")
			}
		} else {
			w.WriteHeader(http.StatusNotImplemented)
			utils.SendMessage(w, "Invalid URL or request")
		}
	}
}
