package store

import (
	"context"
	"database/sql"
	"log"
)

type Medicine struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	MainComponent string  `json:"maincomponent"`
	Price         float32 `json:"price"`
}

type MedicinesStore struct {
	db *sql.DB
}

func (s *MedicinesStore) Create(ctx context.Context, m *Medicine) error {
	query := `INSERT INTO medicamentos(nombre,componenteprincipal,precio)
          VALUES(?,?,?) RETURNING id`

	var id int
	err := s.db.QueryRowContext(ctx, query, m.Name, m.MainComponent, m.Price).Scan(&id)
	if err != nil {
		return err
	}
	log.Printf("Query completed for new medicine with id: %v", id)
	return nil
}

func (s *MedicinesStore) GetByID(ctx context.Context, id string) (*Medicine, error) {
	query := `SELECT * FROM medicamentos WHERE id=?`

	m := Medicine{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&m.Id, &m.Name, &m.MainComponent, &m.Price)
	if err != nil {
		return nil, err
	}
	log.Printf("Query completed for the requested item\n%+v", m)
	return &m, nil
}

func (s *MedicinesStore) GetAll(ctx context.Context) (*[]Medicine, error) {
	query := `SELECT * FROM medicamentos`
	var medicines []Medicine
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := Medicine{}
		err := rows.Scan(&item.Id, &item.Name, &item.MainComponent, &item.Price)
		if err != nil {
			log.Println("Error")
			return &medicines, err
		}
		log.Printf("storing item: %+v", item)
		medicines = append(medicines, item)
	}
	return &medicines, nil
}

func (s *MedicinesStore) GetPaginated(ctx context.Context, limit int, offset int) (*[]Medicine, error) {
	query := `SELECT * FROM medicamentos LIMIT ? OFFSET ?`
	var medicines []Medicine
	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		log.Println("Error")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := Medicine{}
		err := rows.Scan(&item.Id, &item.Name, &item.MainComponent, &item.Price)
		if err != nil {
			log.Println("Error")
			return &medicines, err
		}
		log.Printf("storing item: %+v", item)
		medicines = append(medicines, item)
	}
	return &medicines, nil
}
