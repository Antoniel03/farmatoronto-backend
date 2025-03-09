package main

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

// Para emular la BD
type application struct {
	config config
	store  store.Storage
}
type config struct {
	addr    string
	db      dbConfig
	jwtAuth jwtConfig
}

type dbConfig struct {
	addr string
}

type jwtConfig struct {
	tokenAuth  jwtauth.JWTAuth
	expiration int64
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      enableCORS(mux),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	return srv.ListenAndServe()
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	tkAuth := &app.config.jwtAuth.tokenAuth

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		//TODO retornar datos acorde a la vista
		r.Route("/medicines", func(r chi.Router) {
			// r.Get("/catalog", app.getCatalogHandler)
			r.Get("/", app.getMedicinesHandler)
			r.Get("/{id}", app.getMedicineHandler)
			r.Post("/", app.createMedicineHandler)
		})

		r.Route("/adminview", func(r chi.Router) {
			r.Get("/medicines", app.getMedicinesViewHandler)
			r.Get("/employees", app.getEmployeesViewHandler)
		})

		r.Route("/drugs", func(r chi.Router) {
			r.Post("/", app.CreateDrugsHandler)
		})

		r.Route("/employees", func(r chi.Router) {
			r.Post("/", app.createEmployeeHandler)
			r.Get("/", app.getEmployeesHandler)
			r.Get("/{id}", app.getEmployeeHandler)
		})

		r.Post("/register", app.createUserHandler)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", app.loginHandler)
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(tkAuth))
				r.Use(jwtauth.Authenticator(tkAuth))
				r.Get("/users/{id}", app.getUserHandler)
				r.Get("/users", app.getUsersHandler)
			})
		})
		r.Group(func(r chi.Router) {
			//TODO r.Route lab
			r.Use(jwtauth.Verifier(tkAuth))
			r.Use(jwtauth.Authenticator(tkAuth))
			r.Post("/labs", app.createLabHandler)
			r.Get("/labs/{id}", app.getLabHandler)
			r.Get("/labs", app.getLabsHandler)

			r.Post("/branches", app.createBranchHandler)
			r.Get("/branches/{id}", app.getBranchHandler)
			r.Get("/branches", app.getBranchesHandler)
		})

	})
	return r
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetPaginationParams(query *url.Values) (int, int, error) {
	strLimit := query.Get("limit")
	strOffset := query.Get("offset")

	if strLimit == "" || strOffset == "" {
		return -1, -1, errors.New("invalid params")
	}
	convError := errors.New("conversion error")
	offset, err := strconv.Atoi(strOffset)
	if err != nil {
		return -1, -1, convError
	}

	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		return -1, -1, convError
	}

	return limit, offset, nil
}

func HasPaginationParams(query *url.Values) bool {
	if query.Has("limit") && query.Has("offset") {
		return true
	}
	return false
}

// func GetPaginationParam(query *url.Values) (int, error) {
// 	strPage := query.Get("page")
//
// 	if strPage == "" {
// 		return -1, errors.New("invalid params")
// 	}
// 	convError := errors.New("conversion error")
// 	page, err := strconv.Atoi(strPage)
// 	if err != nil {
// 		return -1, convError
// 	}
// 	return page, nil
// }
