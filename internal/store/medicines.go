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
	query := (`INSERT INTO medicamentos(nombre,componenteprincipal,precio)
          VALUES(?,?,?) RETURNING id`)

	var id int
	err := s.db.QueryRowContext(ctx, query, m.Name, m.MainComponent, m.Price).Scan(&id)
	if err != nil {
		return err
	}
	log.Printf("Query completed for new user with id: %v", id)
	return nil
}

func (s *MedicinesStore) GetByID(ctx context.Context, id int) (*Medicine, error) {
	query := (`SELECT * FROM medicamentos WHERE id=?`)

	m := Medicine{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&m.Id, &m.Name, &m.MainComponent, &m.Price)
	if err != nil {
		return nil, err
	}
	log.Printf("Query completed for the requested item\n%+v", m)
	return &m, nil
}
