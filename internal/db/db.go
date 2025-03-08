package db

import (
	"database/sql"
	"github.com/Antoniel03/farmatoronto-backend/internal/env"
	_ "github.com/mattn/go-sqlite3"

	"os"
)

func SetupDB(database *sql.DB) error {
	// Leer el archivo db_init.sql
	path := env.GetString("DB_SCRIPT", "../../scripts/db_init.sql")
	sqlFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Ejecutar las consultas SQL del archivo
	_, err = database.Exec(string(sqlFile))
	if err != nil {
		return err
	}

	return nil
}

func OpenDB(addr string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", addr)
	if err != nil {
		return nil, err
	}

	return database, nil
}

func CloseDB(db *sql.DB) error {
	return db.Close()
}
