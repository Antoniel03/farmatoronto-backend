package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"github.com/go-chi/chi/v5"
)

func (a *application) createBranchHandler(w http.ResponseWriter, r *http.Request) {
	payload := store.Branch{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	b := &store.Branch{
		CityID:      payload.CityID,
		PhoneNumber: payload.PhoneNumber,
		Address:     payload.Address,
	}
	ctx := r.Context()

	if err := a.store.Branches.Create(ctx, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Couldn't complete operation: ", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *application) getBranchHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	branch, err := a.store.Branches.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(*branch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return

	}
	w.Header().Set("Content-type", "application/json")
}

func (a *application) getBranchesHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	if query.Has("limit") && query.Has("offset") {
		if limit, offset, err := GetPaginationParams(&query); err == nil {
			a.getPaginatedBranches(w, r, limit, offset)
		} else {
			http.Error(w, "invalid page paramameter", http.StatusInternalServerError)
		}
		return
	}

	ctx := r.Context()
	branches, err := a.store.Branches.GetAll(ctx)
	if err != nil {
		http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(branches)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
}

func (a *application) getPaginatedBranches(w http.ResponseWriter, r *http.Request, limit int, offset int) {
	ctx := r.Context()
	branches, err := a.store.Branches.GetPaginated(ctx, limit, offset)
	if err != nil {
		http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(branches)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
}
