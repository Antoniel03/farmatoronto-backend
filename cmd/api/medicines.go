package main

import (
	"encoding/json"
	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

func (a *application) getMedicinesHandler(w http.ResponseWriter, r *http.Request) {
	medicines := []store.Medicine{}
	err := json.NewEncoder(w).Encode(medicines)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
}

func (a *application) getMedicineHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
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
