package main

import (
	"encoding/json"
	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	// "github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func generateHash(password string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Println("Error while hashing: ", err)
		return nil
	}
	return hash
}

func verifyPassword(hash string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		log.Print("Incorrect password")
		return false
	}
	log.Print("Correct password")
	return true
}

func (a *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := store.User{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	hash := generateHash(payload.Password)

	u := &store.User{
		Id:       payload.Id,
		Email:    payload.Email,
		Password: string(hash),
		UserType: payload.UserType,
	}
	ctx := r.Context()

	if err := a.store.Users.Create(ctx, u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Couldn't complete operation: ", err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	log.Printf("Data from the client: %+v", payload)
}

func (a *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	payload := store.User{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("JSON conversion failed:\n", err)
		return
	}
	log.Printf("Received user: %+v", payload)

	ctx := r.Context()

	user, err := a.store.Users.GetByEmail(ctx, payload.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	if !verifyPassword(user.Password, payload.Password) {
		http.Error(w, "Incorrect password", http.StatusBadRequest)
		return
	}
	user.Password = payload.Password

	err = json.NewEncoder(w).Encode(*user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-type", "application/json")
}
