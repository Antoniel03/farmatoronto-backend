package store

import (
	"context"
	"database/sql"
	"log"
)

type Drug struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DrugsStore struct {
	db *sql.DB
}

func (s *DrugsStore) Create(ctx context.Context, d *Drug) error {
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

func (s *DrugsStore) GetAll(ctx context.Context) (*[]Drug, error) {
	query := `SELECT * FROM monodrogas`
	var drugs []Drug
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := Drug{}
		err := rows.Scan(&item.ID, &item.Name, &item.Description)
		if err != nil {
			log.Println(err)
			return &drugs, err
		}
		log.Printf("storing item: %+v", item)
		drugs = append(drugs, item)
	}
	return &drugs, nil
}

func (s *DrugsStore) GetPaginated(ctx context.Context, limit int, offset int) (*[]Drug, bool, error) {
	query := `SELECT * FROM monodrogas LIMIT ? OFFSET ?`
	var drugs []Drug
	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, false, err
	}
	defer rows.Close()

	for rows.Next() {
		item := Drug{}
		err := rows.Scan(&item.ID, &item.Name, &item.Description)
		if err != nil {
			log.Println(err)
			return &drugs, false, err
		}
		log.Printf("storing item: %+v", item)
		drugs = append(drugs, item)
	}

	hasNextPage := handleDrugPagination(s.db, ctx, limit+offset)
	return &drugs, hasNextPage, nil
}

func handleDrugPagination(db *sql.DB, ctx context.Context, nextOffset int) bool {
	hasNextPage := false
	nextQuery := `SELECT COUNT(*) FROM monodrogas`
	row := db.QueryRowContext(ctx, nextQuery)
	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Println("Error getting count: ", err)
	}

	if nextOffset < count {
		hasNextPage = true
	}
	return hasNextPage
}
