package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func SetupDB(database *sql.DB) error {
	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS medicamentos(
			id INTEGER NOT NULL PRIMARY KEY,
			nombre TEXT,
			componenteprincipal TEXT,
			precio TEXT
		)
	`)

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
