package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Antoniel03/farmatoronto-backend/internal/store"
)

func (a *application) createActionHandler(w http.ResponseWriter, r *http.Request) {
	payload := store.Action{}
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		log.Println("test", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	action := payload.Description
	ctx := r.Context()

	if err := a.store.Actions.Create(ctx, action); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Couldn't complete operation: ", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *application) getActionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	actions, err := a.store.Actions.GetAll(ctx)
	if err != nil {
		http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(actions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
}
