package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Antoniel03/farmatoronto-backend/internal/store"
)

func main() {
	obj := store.User{
		Email:    "test4@email.com",
		Password: "plain_text4",
	}

	route := "http://localhost:8082/v1/user/login"
	client := http.Client{}
	data, _ := json.Marshal(obj)

	request, err := http.NewRequest("GET", route, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(strResponse(resp))
	// sendUserData("http://localhost:8082/v1/user", obj)
}

func login(email string, password string) {
}

func sendUserData(route string, user store.User) {
	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Processed data", string(data))
	resp, err := http.Post(route, "application/json", bytes.NewBuffer(data))
	log.Println("Response: " + strResponse(resp))
}

func strResponse(r *http.Response) string {
	str, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(str)
}

func getMedicine(id string) {
	resp, err := http.Get("http://localhost:8082/v1/medicine/" + id)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(strResponse(resp))
}

func health() {
	resp, err := http.Get("http://localhost:8082/v1/health")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(strResponse(resp))
}
