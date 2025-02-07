package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Para emular la BD
var medicamentos = []Medicamento{}

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

func (a *aplicacion) getManejadorMedicamentos(w http.ResponseWriter, r *http.Request) {

	err := json.NewEncoder(w).Encode(medicamentos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
}

func (a *aplicacion) getManejadorMedicamento(w http.ResponseWriter, r *http.Request) {

	_, peticionMedicamento := ExisteMedicamento(chi.URLParam(r, "id"))
	err := json.NewEncoder(w).Encode(peticionMedicamento)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-type", "application/json")
}

func (a *aplicacion) createManejadorMedicamento(w http.ResponseWriter, r *http.Request) {
	var payload Medicamento

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := Medicamento{
		Id:                  payload.Id,
		Nombre:              payload.Nombre,
		ComponentePrincipal: payload.ComponentePrincipal,
		Precio:              payload.Precio,
	}

	if err := InsertarMedicamento(m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	log.Printf("{id:%s, nombre:%s, componente:%s, precio:%s}", payload.Id, payload.Nombre, payload.ComponentePrincipal, payload.Precio)
}

func ExisteMedicamento(id string) (error, Medicamento) {
	for i := range medicamentos {
		if medicamentos[i].Id == id {
			fmt.Println(" -- Medicamento solicitado -- \nID: " + medicamentos[i].Id + "\nNombre: " + medicamentos[i].Nombre)
			return nil, medicamentos[i]
		}
	}

	return errors.New("No se ha encontrado el medicamento"), Medicamento{"", "", "", ""}
}

func InsertarMedicamento(m Medicamento) error {

	if m.Id == "" {
		return errors.New("El id es obligatorio.")
	}

	if m.Nombre == "" {
		return errors.New("El nombre es obligatorio.")
	}

	if m.ComponentePrincipal == "" {
		return errors.New("El componente principal es obligatorio.")
	}

	for _, medicamento := range medicamentos {
		if medicamento.Id == m.Id {
			return errors.New("Ya existe un medicamento con esa id")
		}
	}

	medicamentos = append(medicamentos, m)
	return nil
}
