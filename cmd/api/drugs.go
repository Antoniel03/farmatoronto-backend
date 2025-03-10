package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Antoniel03/farmatoronto-backend/internal/store"
)

func (a *application) createDrugsHandler(w http.ResponseWriter, r *http.Request) {
	payload := store.Drug{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	d := &store.Drug{
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

func (a *application) getDrugsViewHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if !query.Has("limit") || !query.Has("offset") {
		http.Error(w, "invalid pagination paramameters", http.StatusBadRequest)
		return
	}

	if limit, offset, err := GetPaginationParams(&query); err == nil {
		ctx := r.Context()
		drugs, hasNextPage, err := a.store.Drugs.GetPaginated(ctx, limit, offset)
		if err != nil {
			http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(map[string]interface{}{"items": drugs, "nextpage": hasNextPage})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/json")

	} else {
		http.Error(w, "invalid page paramameter", http.StatusInternalServerError)
		return
	}

}

func (a *application) getDrugsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	drugs, err := a.store.Drugs.GetAll(ctx)
	if err != nil {
		http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(drugs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
}
