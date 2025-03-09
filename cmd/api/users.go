package main

import (
	"encoding/json"
	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"github.com/go-chi/chi/v5"

	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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
		ID:         payload.ID,
		EmployeeID: payload.EmployeeID,
		Email:      payload.Email,
		Password:   string(hash),
	}
	ctx := r.Context()

	if err := a.store.Users.Create(ctx, u); err != nil {
		if err.Error() == "UNIQUE constraint failed: usuarios.correo" {
			http.Error(w, "The email used email already exists", http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)

	log.Printf("Data from the client: %+v", payload)
}

func (a *application) getUsersHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	if query.Has("limit") && query.Has("offset") {
		if limit, offset, err := GetPaginationParams(&query); err == nil {
			a.getPaginatedUsers(w, r, limit, offset)
		} else {
			http.Error(w, "invalid page paramameter", http.StatusBadRequest)
		}
		return
	}

	ctx := r.Context()
	users, err := a.store.Users.GetAll(ctx)
	if err != nil {
		http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
}

func (a *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	payload := store.User{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("JSON conversion failed:\n", err)
		return
	}
	log.Printf("Received user: %+v", payload)

	ctx := r.Context()

	user, employee, err := a.store.Users.GetLoginData(ctx, payload.Email)
	if err != nil {
		http.Error(w, "Email not found", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if !verifyPassword(user.Password, payload.Password) {
		http.Error(w, "Incorrect password", http.StatusBadRequest)
		return
	}
	user.Password = payload.Password
	token, err := GenerateJWT(user, employee, a.config.jwtAuth.expiration, &a.config.jwtAuth.tokenAuth)
	if err != nil {
		log.Println("Token error: ", err)
	}
	log.Println("jwt: ", token)

	err = json.NewEncoder(w).Encode(map[string]string{"token": token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-type", "application/json")
}

func (a *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "URL parameter not found", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	user, err := a.store.Users.GetByID(ctx, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Error encoding to JSON", http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/json")

}

func (a *application) getPaginatedUsers(w http.ResponseWriter, r *http.Request, limit int, offset int) {
	ctx := r.Context()
	users, err := a.store.Users.GetPaginated(ctx, limit, offset)
	if err != nil {
		http.Error(w, "Error while retrieveng items", http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
}
