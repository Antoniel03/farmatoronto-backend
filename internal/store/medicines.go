package store

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/Antoniel03/farmatoronto-backend/internal/env"
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

type MedicineView struct {
	Medicine
	Amount  int    `json:"amount"`
	LabName string `json:"lab_name"`
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

func (s *MedicinesStore) GetFiltered(ctx context.Context, limit int, offset int, branch string, drugSubstance string) (*[]MedicineView, error) {
	sql, err := os.ReadFile(env.GetString("MED_Q", "../../internal/store/querys/medicines_view.sql"))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	params, args := handleMedicineFilters(branch, drugSubstance, limit, offset)
	query := string(sql) + params
	log.Println(query)
	var medicines []MedicineView
	rows, err := s.db.QueryContext(ctx, query, *args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := MedicineView{}
		err := rows.Scan(&item.ID, &item.Name, &item.Presentation,
			&item.MainComponent, &item.LabName, &item.Price, &item.Amount, &item.Action)
		if err != nil {
			log.Println(err)
			return &medicines, err
		}
		log.Printf("storing item: %+v", item)
		medicines = append(medicines, item)
	}
	return &medicines, nil
}

func handleMedicineFilters(branch string, drugSubstance string, limit int, offset int) (string, *[]interface{}) {
	var args []interface{}
	finalQuery := ""

	if drugSubstance != "" {
		finalQuery = `JOIN Medic_monodrogas ON Medic_monodrogas.codmedicamento = medicamentos.id
                  JOIN monodrogas ON monodrogas.id = Medic_monodrogas.codmonodroga
                  WHERE monodrogas.nombre=?`
		args = []interface{}{drugSubstance}
		if branch != "" {
			finalQuery += " AND farmacia_sucursal.nombre=?"
			args = append(args, branch)
		}
	} else if branch != "" {
		finalQuery = " WHERE farmacia_sucursal.nombre=?"
		args = []interface{}{branch}
	}
	log.Println(len(args))
	args = append(args, limit)
	args = append(args, offset)
	if len(args) == 0 {
		args = []interface{}{limit, offset}
	}
	finalQuery += " LIMIT ? OFFSET ?"
	return finalQuery, &args
}
