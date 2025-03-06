package store

import (
	"context"
	"database/sql"
	"log"
)

type Branch struct {
	Id      string `json:"id"`
	CityID  int64  `json:"city_id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	// Email   string `json:"email"`
}

type BranchesStore struct {
	db *sql.DB
}

func (s *BranchesStore) Create(ctx context.Context, b *Branch) error {
	query := `INSERT INTO farmacia_sucursal(ciudad_id,nombre,direccion)
          VALUES(?,?,?) RETURNING id`

	var id int64
	err := s.db.QueryRowContext(ctx, query, b.CityID, b.Name, b.Address).Scan(&id)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Query completed for new medicine with id: %v", id)
	return nil
}

func (s *BranchesStore) GetByID(ctx context.Context, id string) (*Branch, error) {
	query := `SELECT * FROM farmacia_sucursal WHERE id=?`

	b := Branch{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&b.Id, &b.CityID, &b.Name, &b.Address)
	if err != nil {
		return nil, err
	}
	log.Printf("Query completed for the requested item\n%+v", b)
	return &b, nil
}

func (s *BranchesStore) GetAll(ctx context.Context) (*[]Branch, error) {
	query := `SELECT * FROM farmacia_sucursal`
	var branches []Branch
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := Branch{}
		err := rows.Scan(&item.Id, &item.CityID, &item.Name, &item.Address)
		if err != nil {
			log.Println("Error")
			return &branches, err
		}
		log.Printf("storing item: %+v", item)
		branches = append(branches, item)
	}
	return &branches, nil
}
