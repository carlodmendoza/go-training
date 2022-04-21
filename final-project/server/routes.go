package main

import (
	"net/http"
	"server/internal/auth"
	"server/internal/categories"
	"server/internal/transactions"
	"server/storage"

	"github.com/go-chi/chi/v5"
)

func GetRouter(db storage.Service) *chi.Mux {
	r := chi.NewRouter()

	// routes not needing authentication
	r.Group(func(r chi.Router) {
		r.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
			auth.Signin(db, w, r)
		})
		r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
			auth.Signup(db, w, r)
		})
	})

	// routes requiring a valid server-generated token
	r.Group(func(r chi.Router) {
		r.Use(auth.Authenticator(db))

		r.Get("/categories", func(w http.ResponseWriter, r *http.Request) {
			categories.ProcessCategories(db, w, r)
		})

		r.Get("/transactions", func(w http.ResponseWriter, r *http.Request) {
			transactions.ListHandler(db, w, r)
		})
		r.Post("/transactions", func(w http.ResponseWriter, r *http.Request) {
			transactions.CreateHandler(db, w, r)
		})
		r.Delete("/transactions", func(w http.ResponseWriter, r *http.Request) {
			transactions.Purge(db, w, r)
		})

		r.Get("/transactions/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
			transactions.RetrieveHandler(db, w, r)
		})
		r.Put("/transactions/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
			transactions.UpdateHandler(db, w, r)
		})
		r.Delete("/transactions/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
			transactions.DeleteHandler(db, w, r)
		})
	})

	return r
}
