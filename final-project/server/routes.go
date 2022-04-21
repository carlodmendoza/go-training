package main

import (
	"net/http"
	"server/internal/auth"
	"server/internal/categories"
	"server/storage"

	"github.com/go-chi/chi/v5"
)

func GetRouter(db storage.Service) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
			auth.Signin(db, w, r)
		})
		r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
			auth.Signup(db, w, r)
		})
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.Authenticator(db))

		r.Get("/categories", func(w http.ResponseWriter, r *http.Request) {
			categories.ProcessCategories(db, w, r)
		})
	})

	return r
}
