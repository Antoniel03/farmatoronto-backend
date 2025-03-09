package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"github.com/go-chi/chi/v5"
)

func (a *application) createLabHandler(w http.ResponseWriter, r *http.Request) {

	payload := store.Lab{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	l := &store.Lab{
		Name:        payload.Name,
		Address:     payload.Address,
		PhoneNumber: payload.PhoneNumber,
	}

	log.Printf("Received: %+v", l)
	ctx := r.Context()

	if err := a.store.Labs.Create(ctx, l); err != nil {
		log.Println("Couldn't complete operation: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *application) getLabHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	lab, err := a.store.Labs.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "lab not found", http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(*lab)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return

	}
	w.Header().Set("Content-type", "application/json")
}

func (a *application) getLabsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	if query.Has("limit") && query.Has("offset") {
		if limit, offset, err := GetPaginationParams(&query); err == nil {
			a.getPaginatedLabs(w, r, limit, offset)
		} else {
			http.Error(w, "invalid page paramameter", http.StatusInternalServerError)
		}
		return
	}

	ctx := r.Context()
	labs, err := a.store.Labs.GetAll(ctx)
	if err != nil {
		http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(labs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
}

func (a *application) getPaginatedLabs(w http.ResponseWriter, r *http.Request, limit int, offset int) {
	ctx := r.Context()
	labs, hasNextPage, err := a.store.Labs.GetPaginated(ctx, limit, offset)
	if err != nil {
		http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(map[string]interface{}{"items": labs, "nextpage": hasNextPage})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")

}
