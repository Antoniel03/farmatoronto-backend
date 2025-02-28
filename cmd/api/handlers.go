package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Antoniel03/farmatoronto-backend/internal/models"
	"github.com/go-chi/chi/v5"
)

func (a *aplicacion) getManejadorMedicamentos(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(medicamentos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
}

func (a *aplicacion) getManejadorMedicamento(w http.ResponseWriter, r *http.Request) {
	_, peticionMedicamento := ExisteMedicamento(chi.URLParam(r, "id"))
	err := json.NewEncoder(w).Encode(peticionMedicamento)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-type", "application/json")
}

func (a *aplicacion) createManejadorMedicamento(w http.ResponseWriter, r *http.Request) {
	var payload models.Medicamento
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := models.Medicamento{
		Id:                  payload.Id,
		Nombre:              payload.Nombre,
		ComponentePrincipal: payload.ComponentePrincipal,
		Precio:              payload.Precio,
	}

	if err := InsertarMedicamento(m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	log.Printf("{id:%s, nombre:%s, componente:%s, precio:%s}", payload.Id, payload.Nombre, payload.ComponentePrincipal, payload.Precio)
}

//TODO move non handler functions out later

func ExisteMedicamento(id string) (error, models.Medicamento) {
	for i := range medicamentos {
		if medicamentos[i].Id == id {
			fmt.Println(" -- Medicamento solicitado -- \nID: " + medicamentos[i].Id + "\nNombre: " + medicamentos[i].Nombre)
			return nil, medicamentos[i]
		}
	}

	return errors.New("No se ha encontrado el medicamento"), models.Medicamento{Id: "", Nombre: "", ComponentePrincipal: "", Precio: ""}
}

func InsertarMedicamento(m models.Medicamento) error {

	if m.Id == "" {
		return errors.New("El id es obligatorio.")
	}

	if m.Nombre == "" {
		return errors.New("El nombre es obligatorio.")
	}

	if m.ComponentePrincipal == "" {
		return errors.New("El componente principal es obligatorio.")
	}

	for _, medicamento := range medicamentos {
		if medicamento.Id == m.Id {
			return errors.New("Ya existe un medicamento con esa id")
		}
	}

	medicamentos = append(medicamentos, m)
	return nil
}
