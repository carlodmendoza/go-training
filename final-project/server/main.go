package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/carlodmendoza/go-training/final-project/server/storage/filebased"
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
