package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	// "github.com/go-chi/chi/v5"
)

func (a *application) CreateDrugsHandler(w http.ResponseWriter, r *http.Request) {
	payload := store.Drugs{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	d := &store.Drugs{
		Name:        payload.Name,
		Description: payload.Description,
	}
	ctx := r.Context()

	if err := a.store.Drugs.Create(ctx, d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Couldn't complete operation: ", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
