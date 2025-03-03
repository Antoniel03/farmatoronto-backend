package main

import (
	"encoding/json"
	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"log"
	"net/http"
)

func (a *application) getEmployeeHandler(w http.ResponseWriter, r *http.Request) {
}

func (a *application) getEmployeesHandler(w http.ResponseWriter, r *http.Request) {
}

func (a *application) createEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var payload store.Employee
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	employee := &store.Employee{
		Id:          payload.Id,
		Name:        payload.Name,
		Lastname:    payload.Lastname,
		Birthday:    payload.Birthday,
		Direction:   payload.Direction,
		PhoneNumber: payload.PhoneNumber,
		Email:       payload.Email,
	}

	//TODO query to db

	ctx := r.Context()
	if err := a.store.Employees.Create(ctx, employee); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
