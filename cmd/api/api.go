package main

import (
	"net/http"
	"time"

	"github.com/Antoniel03/farmatoronto-backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Para emular la BD
var medicamentos = []models.Medicamento{}

type aplicacion struct {
	config config
}

type config struct {
	addr string
}

func (app *aplicacion) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      enableCORS(mux),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	return srv.ListenAndServe()
}

func (app *aplicacion) mount() http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Get("/medicamentos", app.getManejadorMedicamentos)

		r.Route("/medicamento", func(r chi.Router) {
			r.Post("/", app.createManejadorMedicamento)
			r.Get("/{id}", app.getManejadorMedicamento)
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

//Manejadores
