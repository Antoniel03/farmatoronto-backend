package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func setupDB() error {
	// Leer el archivo db_init.sql
	sqlFile, err := os.ReadFile("db_init.sql")
	if err != nil {
		return err
	}

	// Ejecutar las consultas SQL del archivo
	_, err = DB.Exec(string(sqlFile))
	if err != nil {
		return err
	}

	return nil
}

func openDB() error {
	db, err := sql.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		log.Println("Error al abrir la base de datos:", err)
		return err
	}

	DB = db
	return nil
}

func closeDB() error {
	return DB.Close()
}
