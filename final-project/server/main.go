package main

import (
	"final-project/server/models"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Server running in port 8080")
	db := &models.Database{}
	db.InitializeDatabase()
	if err := http.ListenAndServe("localhost:8080", db.Handler()); err != nil {
		log.Fatalf("Error ListenAndServe(): %s", err.Error())
	}
}
