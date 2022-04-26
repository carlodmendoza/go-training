package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/carlodmendoza/go-training/final-project/server/storage/filebased"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Debug().Msg("Server running on port 8080")

	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) //nolint

	// TODO: create env var for file path and chosen storage service
	storage := filebased.Initialize("../deploy/dev/server/storage/data")
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

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}
}
