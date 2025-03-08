package store

import (
	"context"
	"database/sql"
	"log"
)

type MedicineExtraData struct {
	mlID       int64 `json:"ml_id"`
	MedicineID int64 `json:"id_medicine"`
	BranchID   int64 `json:"branch_id"`
	LabID      int64 `json:"branch_id"`
	Amount     int   `json:"amount"`
}

type Medicine struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Presentation  string  `json:"presentation"`
	MainComponent string  `json:"maincomponent"`
	Action        string  `json:"action"`
	Price         float32 `json:"price"`
}

type MedicinesStore struct {
	db *sql.DB
}

func (s *MedicinesStore) Create(ctx context.Context, m *Medicine, extraData *MedicineExtraData) error {
	query := `INSERT INTO medicamentos(nombre,componenteprincipal,presentacion,accion,precio)
          VALUES(?,?,?,?,?) RETURNING id`

	var medicineID int
	err := s.db.QueryRowContext(ctx, query, m.Name, m.MainComponent, m.Presentation, m.Action, m.Price).Scan(&medicineID)
	if err != nil {
		return err
	}

	query = `INSERT INTO Medic_labs(codmedicamento,codlaboratorio)
          VALUES(?,?) RETURNING id`

	var medicLabsID int
	err = s.db.QueryRowContext(ctx, query, medicineID, extraData.LabID).Scan(&medicLabsID)
	if err != nil {
		return err
	}

	query = `INSERT INTO stock(farmacia_id,medicamento_id,cantidad)
          VALUES(?,?,?) RETURNING id`

	var stockID int
	err = s.db.QueryRowContext(ctx, query, extraData.BranchID, medicineID, extraData.Amount).Scan(&stockID)
	if err != nil {
		return err
	}
	return nil
}

func (s *MedicinesStore) GetByID(ctx context.Context, id string) (*Medicine, error) {
	query := `SELECT * FROM medicamentos WHERE id=?`

	m := Medicine{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&m.ID, &m.Name, &m.Presentation, &m.MainComponent, &m.Action, &m.Price)
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
		err := rows.Scan(&item.ID, &item.Name, &item.Presentation, &item.MainComponent, &item.Action, &item.Price)
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
		err := rows.Scan(&item.ID, &item.Name, &item.Presentation, &item.MainComponent, &item.Action, &item.Price)
		if err != nil {
			log.Println(err)
			return &medicines, err
		}
		log.Printf("storing item: %+v", item)
		medicines = append(medicines, item)
	}
	return &medicines, nil
}
