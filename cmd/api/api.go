package main

import (
	"net/http"
	"time"

	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Para emular la BD
type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	addr string
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

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Get("/medicines", app.getMedicinesHandler)
		// r.Route("/user", func(r chi.Router)

		r.Route("/medicine", func(r chi.Router) {
			r.Post("/create", app.createMedicineHandler)
			r.Get("/{id}", app.getMedicineHandler)
			r.Put("/{id}", app.createMedicineHandler)
			r.Delete("/", app.createMedicineHandler)
		})

		r.Route("/employee", func(r chi.Router) {
			r.Post("/", app.createEmployeeHandler)
			r.Get("/", app.getEmployeesHandler)
			r.Get("/{id}", app.getEmployeeHandler)
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/login", app.loginHandler)
			r.Post("/create", app.createUserHandler)
			// r.Get("/{email}", app.getUserHandler)
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
