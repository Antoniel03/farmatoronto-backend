package main

import (
	"net/http"
	"time"

	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

// Para emular la BD
type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr    string
	db      dbConfig
	jwtAuth jwtConfig
}

type dbConfig struct {
	addr string
}

type jwtConfig struct {
	tokenAuth  jwtauth.JWTAuth
	expiration int64
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      enableCORS(mux),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	return srv.ListenAndServe()
}

//TODO middleware de autorizacion

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	tkAuth := &app.config.jwtAuth.tokenAuth

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		// r.Route("/user", func(r chi.Router)

		r.Route("/medicines", func(r chi.Router) {
			r.Get("/", app.getMedicinesHandler)
			r.Get("/{id}", app.getMedicineHandler)
			r.Post("/", app.createMedicineHandler)
			// r.Put("/{id}", app.createMedicineHandler)
			// r.Delete("/{id}", app.createMedicineHandler)
		})

		r.Route("/employees", func(r chi.Router) {
			r.Post("/", app.createEmployeeHandler)
			r.Get("/", app.getEmployeesHandler)
			r.Get("/{id}", app.getEmployeeHandler)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", app.loginHandler)
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(tkAuth))
				r.Use(jwtauth.Authenticator(tkAuth))
				r.Post("/register", app.createUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(tkAuth))
				r.Use(jwtauth.Authenticator(tkAuth))
				r.Get("/users/{id}", app.getUserHandler)
				r.Get("/users", app.getUsersHandler)

			})
		})

	})
	return r
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
