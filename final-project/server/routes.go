package main

import (
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

		r.Method("GET", "/categories", http.Handler{Storage: db, Function: categories.ProcessCategories})

		r.Method("GET", "/transactions", http.Handler{Storage: db, Function: transactions.ListHandler})
		r.Method("POST", "/transactions", http.Handler{Storage: db, Function: transactions.CreateHandler})
		r.Method("DELETE", "/transactions", http.Handler{Storage: db, Function: transactions.Purge})

		r.Method("GET", "/transactions/{id:[0-9]+}", http.Handler{Storage: db, Function: transactions.RetrieveHandler})
		r.Method("PUT", "/transactions/{id:[0-9]+}", http.Handler{Storage: db, Function: transactions.UpdateHandler})
		r.Method("DELETE", "/transactions/{id:[0-9]+}", http.Handler{Storage: db, Function: transactions.DeleteHandler})
	})

	return r
}
