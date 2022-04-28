package http

import (
	"net/http"

	"github.com/carlodmendoza/go-training/final-project/server/storage"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	Storage  storage.Service
	Function func(storage.Service, *ResponseWriter, *http.Request) error
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := &ResponseWriter{w: w}

	err := h.Function(h.Storage, rw, r)
	if err != nil {
		switch e := err.(type) {
		case HandlerError:
			rw.WriteHeader(e.Status())
			rw.WriteError(e.Error()) //nolint
			log.Error().Err(err).Msgf("HTTP %d: %s\n", e.Status(), e)
		}
	}
}
