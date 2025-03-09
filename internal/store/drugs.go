package store

import (
	"context"
	"database/sql"
	"log"
)

type Drugs struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DrugsStore struct {
	db *sql.DB
}

func (s *DrugsStore) Create(ctx context.Context, d *Drugs) error {
	query := `INSERT INTO monodrogas(nombre,descripcion)
          VALUES(?,?) RETURNING id`

	var id int64
	err := s.db.QueryRowContext(ctx, query, d.Name, d.Description).Scan(&id)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Query completed for new drug with id: %v", id)
	return nil
}
