package store

import (
	"context"
	"database/sql"
	"log"
)

type Branch struct {
	ID          int64  `json:"id"`
	CityID      int64  `json:"city_id"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phonenumber"`
}

type BranchesStore struct {
	db *sql.DB
}

func (s *BranchesStore) Create(ctx context.Context, b *Branch) error {
	query := `INSERT INTO farmacia_sucursal(ciudad_id,telefono,direccion)
          VALUES(?,?,?) RETURNING id`

	var id int64
	err := s.db.QueryRowContext(ctx, query, b.CityID, b.PhoneNumber, b.Address).Scan(&id)
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
	err := s.db.QueryRowContext(ctx, query, id).Scan(&b.ID, &b.CityID, &b.Address, &b.PhoneNumber)
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
		err := rows.Scan(&item.ID, &item.CityID, &item.Address, &item.PhoneNumber)
		if err != nil {
			log.Println("Error")
			return &branches, err
		}
		log.Printf("storing item: %+v", item)
		branches = append(branches, item)
	}
	return &branches, nil
}

func (s *BranchesStore) GetPaginated(ctx context.Context, limit int, offset int) (*[]Branch, bool, error) {
	query := `SELECT * FROM farmacia_sucursal LIMIT ? OFFSET ?`
	var branches []Branch
	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, false, err
	}
	defer rows.Close()

	for rows.Next() {
		item := Branch{}
		err = rows.Scan(&item.ID, &item.CityID, &item.Address, &item.PhoneNumber)
		if err != nil {
			log.Println(err)
			return &branches, false, err
		}
		log.Printf("storing item: %+v", item)
		branches = append(branches, item)
	}

	hasNextPage := false
	nextQuery := `SELECT COUNT(*) FROM farmacia_sucursal`
	nextOffset := limit + offset
	row := s.db.QueryRowContext(ctx, nextQuery)
	var count int
	err = row.Scan(&count)
	if err != nil {
		log.Println("Error getting count: ", err)
	}

	if nextOffset < count {
		hasNextPage = true

	}
	return &branches, hasNextPage, nil
}
