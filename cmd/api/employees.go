package main

import (
	"encoding/json"
	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func (a *application) getEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	employee, err := a.store.Employees.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(*employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-type", "application/json")
}

func (a *application) getEmployeesHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	employees, err := a.store.Employees.GetAll(ctx)
	if err != nil {
		http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(employees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
}

func (a *application) createEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var payload store.Employee
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	employee := &store.Employee{
		ID:          payload.ID,
		Name:        payload.Name,
		Lastname:    payload.Lastname,
		Birthday:    payload.Birthday,
		Address:     payload.Address,
		PhoneNumber: payload.PhoneNumber,
	}

	ctx := r.Context()
	if err := a.store.Employees.Create(ctx, employee); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *application) getEmployeesViewHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit, offset, err := GetPaginationParams(&query)
	if err != nil {
		http.Error(w, "invalid parameters", http.StatusBadRequest)
		return
	}
	branch := query.Get("branch")

	log.Println("branch: " + branch)
	ctx := r.Context()
	employees, err := a.store.Employees.GetFiltered(ctx, limit, offset, branch)
	if err != nil {
		http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(employees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")

}
