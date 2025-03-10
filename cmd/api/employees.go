package main

import (
	"encoding/json"
	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type EmployeeData struct {
	store.Employee
	store.User
	BranchID int64  `json:"branch_id"`
	Date     string `json:"date"`
}

func (a *application) registerEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var payload EmployeeData
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash := generateHash(payload.Password)

	user := &store.User{
		Email:    payload.Email,
		Password: string(hash),
	}

	employee := &store.Employee{
		Name:        payload.Name,
		Lastname:    payload.Lastname,
		Birthday:    payload.Birthday,
		Address:     payload.Address,
		PhoneNumber: payload.PhoneNumber,
		Role:        payload.Role,
		C_ID:        payload.C_ID,
	}

	ctx := r.Context()
	if err := a.store.Employees.RegisterEmployee(ctx, employee, user, payload.BranchID, payload.Date); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

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
		Name:        payload.Name,
		Lastname:    payload.Lastname,
		Birthday:    payload.Birthday,
		Address:     payload.Address,
		PhoneNumber: payload.PhoneNumber,
		Role:        payload.Role,
		C_ID:        payload.C_ID,
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

	log.Println("branch: "+branch+"limit: ", limit, "offset: ", offset)
	ctx := r.Context()
	employees, hasNextPage, err := a.store.Employees.GetFiltered(ctx, limit, offset, branch)
	if err != nil {
		http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(map[string]interface{}{"items": employees, "nextpage": hasNextPage})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")

}
