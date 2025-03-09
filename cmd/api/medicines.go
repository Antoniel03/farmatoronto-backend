package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"github.com/go-chi/chi/v5"
)

type MedicinePayload struct {
	store.Medicine
	ActionDescription string `json:"action_description"`
	LabID             int64  `json:"lab_id"`
	BranchID          int64  `json:"branch_id"`
	Amount            int    `json:"amount"`
}

func (a *application) getMedicinesHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	if query.Has("limit") && query.Has("offset") {
		if limit, offset, err := GetPaginationParams(&query); err == nil {
			a.getPaginatedMedicines(w, r, limit, offset)
		} else {
			http.Error(w, "invalid page paramameter", http.StatusBadRequest)
		}
		return
	}

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
	payload := MedicinePayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	m := &store.Medicine{
		Name:          payload.Name,
		Presentation:  payload.Presentation,
		ActionID:      payload.ActionID,
		MainComponent: payload.MainComponent,
		Price:         payload.Price,
	}
	ctx := r.Context()

	extra := &store.MedicineExtraData{
		LabID:    payload.LabID,
		BranchID: payload.BranchID,
		Amount:   payload.Amount,
	}

	if err := a.store.Medicines.Create(ctx, m, extra); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Couldn't complete operation: ", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *application) getPaginatedMedicines(w http.ResponseWriter, r *http.Request, limit int, offset int) {
	ctx := r.Context()
	medicines, err := a.store.Medicines.GetPaginated(ctx, limit, offset)
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

// func medicineFiltering(r *url.Values) (string, error) {
// 	branch := r.Get("branch")
// 	drugSubstance := r.Get("drug_substance")
//
// 	limit, offset, err := GetPaginationParams(r)
// 	if err != nil {
// 		return "", errors.New("invalid pagination parameters")
// 	}
//
// 	whereClauses := []string{}
// 	args := []interface{}{}
// 	argIndex := 1
//
// 	if branch != "" {
// 		whereClauses = append(whereClauses, "farmacia_sucursal.id=?")
// 		args = append(args, "%"+branch+"%")
// 		argIndex++
// 	}
//
// 	if drugSubstance != "" {
// 		whereClauses = append(whereClauses, "monodrogas.id=?")
// 		args = append(args, "%"+drugSubstance+"%")
// 		argIndex++
// 	}
//
// 	args = append(args, limit, offset)
//
// 	whereSQL := ""
// 	if len(whereClauses) > 0 {
// 		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
// 	}
//
// 	query := fmt.Sprintf("SELECT id, name FROM medicamentos LIMIT $%d OFFSET $%d", whereSQL, argIndex, argIndex+1)
// 	return query, nil
// }

func (a *application) getMedicinesViewHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit, offset, err := GetPaginationParams(&query)
	if err != nil {
		http.Error(w, "invalid parameters", http.StatusBadRequest)
		return
	}
	branch := query.Get("branch")
	drugSubstance := query.Get("drugsubstance")

	log.Println("branch: " + branch + ", drug: " + drugSubstance)
	ctx := r.Context()
	medicines, err := a.store.Medicines.GetFiltered(ctx, limit, offset, branch, drugSubstance)
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

// func medicineFiltering(r *url.Values) (string, error) {
//   GetPaginationParams(r)
//   query:="SELECT medicamentos."
//   if url
//   return "",nil
// }

// func (a *application) getCatalogHandler(w http.ResponseWriter, r *http.Request) {
// 	query := r.URL.Query()
//
// 	if !query.Has("page") {
// 		http.Error(w, "pagination parameters not found", http.StatusBadRequest)
// 		return
// 	}
// 	page, err := GetPaginationParam(&query)
// 	if err != nil {
// 		http.Error(w, "invalid pagination parameter", http.StatusBadRequest)
// 		return
// 	}
// 	limit := 2
// 	a.getPaginatedMedicines(w, r, limit, page*limit)
// }
//
