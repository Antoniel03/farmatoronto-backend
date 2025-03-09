package store

import (
	"context"
	"database/sql"
	"log"
)

type Lab struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phonenumber"`
}

type LabsStore struct {
	db *sql.DB
}

func (s *LabsStore) Create(ctx context.Context, l *Lab) error {
	query := `INSERT INTO laboratorio(nombre,direccion,telefono)
          VALUES(?,?,?) RETURNING id`

	var id int
	err := s.db.QueryRowContext(ctx, query, l.Name, l.Address, l.PhoneNumber).Scan(&id)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Query completed for new medicine with id: %v", id)
	return nil
}

func (s *LabsStore) GetByID(ctx context.Context, id string) (*Lab, error) {
	query := `SELECT * FROM laboratorio WHERE id=?`

	l := Lab{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&l.ID, &l.Name, &l.Address, &l.PhoneNumber)
	if err != nil {
		return nil, err
	}
	log.Printf("Query completed for the requested item\n%+v", l)
	return &l, nil
}

func (s *LabsStore) GetAll(ctx context.Context) (*[]Lab, error) {
	query := `SELECT * FROM laboratorio`
	var labs []Lab
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := Lab{}
		err := rows.Scan(&item.ID, &item.Name, &item.Address, &item.PhoneNumber)
		if err != nil {
			log.Println("Error")
			return &labs, err
		}
		log.Printf("storing item: %+v", item)
		labs = append(labs, item)
	}
	return &labs, nil
}

func (s *LabsStore) GetPaginated(ctx context.Context, limit int, offset int) (*[]Lab, error) {
	query := `SELECT * FROM laboratorio LIMIT ? OFFSET ?`
	var labs []Lab
	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		log.Println("Error")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := Lab{}
		err := rows.Scan(&item.ID, &item.Name, &item.Address, &item.PhoneNumber)
		if err != nil {
			log.Println("Error")
			return &labs, err
		}
		log.Printf("storing item: %+v", item)
		labs = append(labs, item)
	}
	return &labs, nil
}
