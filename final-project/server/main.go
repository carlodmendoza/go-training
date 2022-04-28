package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/carlodmendoza/go-training/final-project/server/storage"
	"github.com/carlodmendoza/go-training/final-project/server/storage/filebased"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().Msg(fmt.Sprintf("Server running on port %v", os.Getenv("HTTP_PORT")))

	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) //nolint

	var storage storage.Service
	switch os.Getenv("STORAGE_SERVICE") {
	case "filebased":
		storage = filebased.Initialize(os.Getenv("FILE_STORAGE_PATH"))
	case "redis":
		// TODO: implement redis
	}
	r := GetRouter(storage)

	go func() {
		<-sigChannel
		err := storage.Shutdown()
		if err != nil {
			log.Error().Err(err).Msg("Shutdown error")
		}
		// TODO: perform graceful shutdown
		log.Fatal().Msg("Shutting down the server")
	}()

	err := http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("HTTP_PORT")), r)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}
}
