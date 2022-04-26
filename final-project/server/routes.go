package main

import (
	gohttp "net/http"
	"time"

	"github.com/carlodmendoza/go-training/final-project/server/internal/auth"
	"github.com/carlodmendoza/go-training/final-project/server/internal/categories"
	"github.com/carlodmendoza/go-training/final-project/server/internal/transactions"
	"github.com/carlodmendoza/go-training/final-project/server/pkg/http"
	"github.com/carlodmendoza/go-training/final-project/server/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func GetRouter(db storage.Service) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// routes not needing authentication
	r.Group(func(r chi.Router) {
		r.Method("POST", "/signin", http.Handler{Storage: db, Function: auth.Signin})
		r.Method("POST", "/signup", http.Handler{Storage: db, Function: auth.Signup})
	})

	// routes requiring a valid server-generated token
	r.Group(func(r chi.Router) {
		r.Use(auth.Authenticator(db))

		r.Get("/categories", func(w gohttp.ResponseWriter, r *gohttp.Request) {
			categories.ProcessCategories(db, w, r)
		})

		r.Get("/transactions", func(w gohttp.ResponseWriter, r *gohttp.Request) {
			transactions.ListHandler(db, w, r)
		})
		r.Post("/transactions", func(w gohttp.ResponseWriter, r *gohttp.Request) {
			transactions.CreateHandler(db, w, r)
		})
		r.Delete("/transactions", func(w gohttp.ResponseWriter, r *gohttp.Request) {
			transactions.Purge(db, w, r)
		})

		r.Get("/transactions/{id:[0-9]+}", func(w gohttp.ResponseWriter, r *gohttp.Request) {
			transactions.RetrieveHandler(db, w, r)
		})
		r.Put("/transactions/{id:[0-9]+}", func(w gohttp.ResponseWriter, r *gohttp.Request) {
			transactions.UpdateHandler(db, w, r)
		})
		r.Delete("/transactions/{id:[0-9]+}", func(w gohttp.ResponseWriter, r *gohttp.Request) {
			transactions.DeleteHandler(db, w, r)
		})
	})

	return r
}
