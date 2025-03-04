package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"github.com/go-chi/chi/v5"
)

func (a *application) getMedicinesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	medicines, err := a.store.Medicines.GetAll(ctx)
	if err != nil {
		http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(medicines)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
}

func (a *application) getMedicineHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	log.Println(err)
	// 	return
	// }
	ctx := r.Context()

	medicine, err := a.store.Medicines.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(*medicine)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return

	}
	w.Header().Set("Content-type", "application/json")
}

func (a *application) createMedicineHandler(w http.ResponseWriter, r *http.Request) {
	payload := store.Medicine{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	m := &store.Medicine{
		Id:            payload.Id,
		Name:          payload.Name,
		MainComponent: payload.MainComponent,
		Price:         payload.Price,
	}
	ctx := r.Context()

	if err := a.store.Medicines.Create(ctx, m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Couldn't complete operation: ", err)
		return
	}

}
