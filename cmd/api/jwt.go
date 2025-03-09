package main

import (
	// "strconv"
	"time"

	"github.com/Antoniel03/farmatoronto-backend/internal/store"
	"github.com/go-chi/jwtauth/v5"
)

func NewJWTAuth(secret string, algorithm string) jwtauth.JWTAuth {
	tokenAuth := jwtauth.New(algorithm, []byte(secret), nil)
	return *tokenAuth
}

func GenerateJWT(u *store.User, e *store.Employee, exp int64, tokenAuth *jwtauth.JWTAuth) (string, error) {
	expiration := time.Second * time.Duration(exp)
	userName := e.Name + " " + e.Lastname
	claims := map[string]interface{}{"id": u.ID, "role": e.Role, "name": userName}
	jwtauth.SetExpiry(claims, time.Now().Add(expiration))
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
