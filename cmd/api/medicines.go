package main

import (
	"encoding/json"
	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

type MedicinePayload struct {
	store.Medicine
	store.MedicineLinks
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
		MainComponent: payload.MainComponent,
		Price:         payload.Price,
		Amount:        payload.Amount,
	}
	ctx := r.Context()

	extra := &store.MedicineLinks{
		LabID:      payload.LabID,
		PharmacyID: payload.PharmacyID,
		ActionID:   payload.ActionID,
		DrugsIDs:   payload.DrugsIDs,
	}

	log.Println(payload.DrugsIDs)
	if err := a.store.Medicines.Create(ctx, m, extra); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Couldn't complete operation: ", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *application) getMedicinesHandler(w http.ResponseWriter, r *http.Request) {
	querys := r.URL.Query()

	if querys.Has("limit") && querys.Has("offset") {
		a.getFilteredMedicinesHandler(w, r)
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
		http.Error(w, "Medicine not found.", http.StatusBadRequest)
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

//

//	func (a *application) getPaginatedMedicines(w http.ResponseWriter, r *http.Request, limit int, offset int) {
//		ctx := r.Context()
//		medicines, err := a.store.Medicines.GetPaginated(ctx, limit, offset)
//		if err != nil {
//			http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
//		}
//		err = json.NewEncoder(w).Encode(medicines)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		w.Header().Set("Content-type", "application/json")
//	}
//
// // func medicineFiltering(r *url.Values) (string, error) {
// // 	pharmacy := r.Get("pharmacy")
// // 	drugSubstance := r.Get("drug_substance")
// //
// // 	limit, offset, err := GetPaginationParams(r)
// // 	if err != nil {
// // 		return "", errors.New("invalid pagination parameters")
// // 	}
// //
// // 	whereClauses := []string{}
// // 	args := []interface{}{}
// // 	argIndex := 1
// //
// // 	if pharmacy != "" {
// // 		whereClauses = append(whereClauses, "farmacia_sucursal.id=?")
// // 		args = append(args, "%"+pharmacy+"%")
// // 		argIndex++
// // 	}
// //
// // 	if drugSubstance != "" {
// // 		whereClauses = append(whereClauses, "monodrogas.id=?")
// // 		args = append(args, "%"+drugSubstance+"%")
// // 		argIndex++
// // 	}
// //
// // 	args = append(args, limit, offset)
// //
// // 	whereSQL := ""
// // 	if len(whereClauses) > 0 {
// // 		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
// // 	}
// //
// // 	query := fmt.Sprintf("SELECT id, name FROM medicamentos LIMIT $%d OFFSET $%d", whereSQL, argIndex, argIndex+1)
// // 	return query, nil
// // }
func (a *application) getFilteredMedicinesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit, offset, err := GetPaginationParams(&query)
	log.Println(limit, offset)
	if err != nil {
		http.Error(w, "invalid parameters", http.StatusBadRequest)
		return
	}
	pharmacy, err := strconv.Atoi(query.Get("pharmacy"))
	if err != nil {
		pharmacy = 0
	}
	drugSubstance, err := strconv.Atoi(query.Get("drug"))
	if err != nil {
		drugSubstance = 0
	}

	log.Println("pharmacy: ", pharmacy, ", drug: ", drugSubstance)
	ctx := r.Context()

	filters := map[string]int{"limit": limit, "offset": offset, "pharmacy": pharmacy, "drug": drugSubstance}
	medicines, hasNextPage, err := a.store.Medicines.GetFiltered(ctx, filters)
	if err != nil {
		http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(map[string]interface{}{"items": medicines, "nextpage": hasNextPage})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
}
