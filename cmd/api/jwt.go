package main

import (
	// "strconv"
	"github.com/go-chi/jwtauth/v5"
	"time"
)

func NewJWTAuth(secret string, algorithm string) jwtauth.JWTAuth {
	tokenAuth := jwtauth.New(algorithm, []byte(secret), nil)
	return *tokenAuth
}

func GenerateJWT(id int, role string, exp int64, tokenAuth *jwtauth.JWTAuth) (string, error) {
	expiration := time.Second * time.Duration(exp)
	claims := map[string]interface{}{"id": id, "role": role}
	jwtauth.SetExpiry(claims, time.Now().Add(expiration))
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
