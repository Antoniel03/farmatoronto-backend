package store

import (
	"context"
	"database/sql"
	"log"
)

type Action struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

type ActionsStore struct {
	db *sql.DB
}

func (s *ActionsStore) Create(ctx context.Context, description string) error {
	query := `INSERT INTO accion_terapeutica(descripcion) VALUES(?) RETURNING id`

	log.Println(description)
	var id int64
	err := s.db.QueryRowContext(ctx, query, description).Scan(&id)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Query completed for new action with id: %v", id)
	return nil
}

func (s *ActionsStore) GetAll(ctx context.Context) (*[]Action, error) {
	query := `SELECT * FROM accion_terapeutica`
	var actions []Action
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := Action{}
		err := rows.Scan(&item.ID, &item.Description)
		if err != nil {
			log.Println(err)
			return &actions, err
		}
		log.Printf("storing item: %+v", item)
		actions = append(actions, item)
	}
	return &actions, nil
}
