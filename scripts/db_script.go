package main

import (
	"log"
	"os"

	"github.com/Antoniel03/farmatoronto-backend/internal/db"
)

func main() {
	path := os.Args[1]

	c, ioErr := os.ReadFile(path)
	if ioErr != nil {
		log.Fatal(ioErr)
	}
	sql := string(c)

	db, err := db.OpenDB("../internal/db/farma_db.db")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec(sql)
}
