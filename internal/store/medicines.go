package store

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/Antoniel03/farmatoronto-backend/internal/env"
)

//Human readable

type Medicine struct {
	ID            int32    `json:"id"`
	Name          string   `json:"name"`
	Presentation  string   `json:"presentation"`
	MainComponent string   `json:"maincomponent"`
	Action        string   `json:"action"`
	Price         float32  `json:"price"`
	Drugs         []string `json:"drugs"`
	LabName       string   `json:"lab_name"`
	PharmacyName  string   `json:"pharmacy_city"`
	Amount        int      `json:"amount"`
}

// For requests
type MedicineLinks struct {
	LabID      int32   `json:"lab_id"`
	PharmacyID int32   `json:"pharmacy_id"`
	ActionID   int32   `json:"action_id"`
	DrugsIDs   []int32 `json:"drugs_ids"`
}

type MedicinesStore struct {
	db *sql.DB
}

// Tables affected: medicamentos, medic_labs, stock, medic_monodrogas
func (s *MedicinesStore) Create(ctx context.Context, m *Medicine, ml *MedicineLinks) error {
	query := `INSERT INTO medicamentos(nombre,componenteprincipal,presentacion,accion_id,precio)
          VALUES(?,?,?,?,?) RETURNING id`

	var medicineID int32
	err := s.db.QueryRowContext(ctx, query, m.Name, m.MainComponent, m.Presentation, ml.ActionID, m.Price).Scan(&medicineID)
	if err != nil {
		return err
	}

	err = linkToLab(s, ctx, medicineID, &ml.LabID)
	if err != nil {
		return err
	}

	err = linkToPharmacy(s, ctx, m, ml)
	if err != nil {
		return err
	}

	err = linkToDrugs(s, ctx, medicineID, ml.DrugsIDs)
	if err != nil {
		return err
	}

	return nil
}

func linkToPharmacy(s *MedicinesStore, ctx context.Context, m *Medicine, ml *MedicineLinks) error {
	query := `INSERT INTO stock(farmacia_id,medicamento_id,cantidad)
          VALUES(?,?,?) RETURNING id`

	var stockID int
	err := s.db.QueryRowContext(ctx, query, ml.PharmacyID, m.ID, m.Amount).Scan(&stockID)
	return err
}

func linkToLab(s *MedicinesStore, ctx context.Context, medicineID int32, labID *int32) error {
	query := `INSERT INTO medic_labs(medicamento_id,laboratorio_id)
          VALUES(?,?) RETURNING id`

	var medicLabsID int
	err := s.db.QueryRowContext(ctx, query, medicineID, labID).Scan(&medicLabsID)
	return err
}

func linkToDrugs(s *MedicinesStore, ctx context.Context, medicineID int32, drugs []int32) error {
	query := `INSERT INTO medic_monodrogas(medicamento_id, monodroga_id) VALUES(?,?) RETURNING medic_monodrogas.id`
	var err error

	log.Println(len(drugs))
	var medicDrugsID int
	for i := 0; i < len(drugs); i++ {
		log.Println("medicine ", medicineID, " linked to ", drugs[i])
		err = s.db.QueryRowContext(ctx, query, medicineID, drugs[i]).Scan(&medicDrugsID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MedicinesStore) GetAll(ctx context.Context) (*[]Medicine, error) {
	query := `SELECT medicamentos.*, accion_terapeutica.descripcion, laboratorio.nombre, ciudad.nombre FROM medicamentos JOIN accion_terapeutica ON accion_terapeutica.id = medicamentos.accion_id
            JOIN medic_labs ON medic_labs.medicamento_id= medicamentos.id
            JOIN laboratorio ON laboratorio.id = medic_labs.laboratorio_id
            JOIN stock ON stock.medicamento_id = medicamentos.id
            JOIN sucursal ON sucursal.id = stock.farmacia_id
            JOIN ciudad ON ciudad.id = sucursal.ciudad_id`
	var medicines []Medicine
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		ignore := ""
		item := Medicine{}
		err := rows.Scan(&item.ID, &item.Name, &item.MainComponent, &ignore, &item.Presentation, &item.Price, &item.Action, &item.LabName, &item.PharmacyName)
		if err != nil {
			log.Println(err)
			return &medicines, err
		}
		log.Printf("storing item: %+v", item)
		drugs, err := GetDrugNames(s.db, ctx, item.ID)
		if err != nil {
			log.Println(err)
			return &medicines, err
		}
		item.Drugs = *drugs
		medicines = append(medicines, item)
	}
	return &medicines, nil
}

func (s *MedicinesStore) GetByID(ctx context.Context, id string) (*Medicine, error) {
	query := `SELECT medicamentos.*, accion_terapeutica.descripcion
             FROM medicamentos JOIN accion_terapeutica 
             ON accion_terapeutica.id = medicamentos.accion_id 
             WHERE id=?`

	m := Medicine{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&m.ID, &m.Name, &m.MainComponent, &m.Presentation, &m.Price, &m.Action)
	if err != nil {
		return nil, err
	}
	log.Printf("Query completed for the requested item\n%+v", m)
	return &m, nil
}

func (s *MedicinesStore) GetFiltered(ctx context.Context, pagination map[string]int) (*[]Medicine, bool, error) {
	path := env.GetString("MED_Q", "internal/store/querys/medicines_view.sql")

	query, err := os.ReadFile(path)
	baseQuery := string(query)
	if err != nil {
		log.Println(err)
		return nil, false, err
	}

	filters, args := handleMedicineFilters(pagination)
	log.Printf("finalQuery: %+v", baseQuery+filters)
	var medicines []Medicine
	rows, err := s.db.QueryContext(ctx, baseQuery+filters, *args...)
	if err != nil {
		log.Println(err)
		return nil, false, err
	}
	defer rows.Close()

	for rows.Next() {
		item := Medicine{}
		err := rows.Scan(&item.ID, &item.Name, &item.MainComponent, &item.Presentation,
			&item.LabName, &item.Price, &item.Amount, &item.Action, &item.PharmacyName)
		if err != nil {
			log.Println(err)
			return &medicines, false, err
		}
		log.Printf("storing item: %+v", item)
		drugs, err := GetDrugNames(s.db, ctx, item.ID)
		if err != nil {
			log.Println(err)
			return &medicines, false, err
		}
		item.Drugs = *drugs
		medicines = append(medicines, item)
	}
	nextPage := handleMedicinePagination(s.db, ctx, pagination)
	return &medicines, nextPage, err
}

//Utilities

func handleMedicineFilters(filters map[string]int) (string, *[]interface{}) {
	var args []interface{}
	finalQuery := ""
	// log.Println("handleMedicinefilters check", limit, offset)

	if filters["drug"] != 0 {
		finalQuery += ` JOIN medic_monodrogas ON medic_monodrogas.medicamento_id = medicamentos.id
                    JOIN monodrogas ON monodrogas.id = medic_monodrogas.monodroga_id`
	}

	if filters["pharmacy"] != 0 {
		finalQuery += `JOIN ciudad ON ciudad.id = sucursal.ciudad_id`
	}

	if filters["drug"] != 0 {
		finalQuery += ` WHERE monodrogas.id=?`
		args = []interface{}{filters["drug"]}
		if filters["pharmacy"] != 0 {
			finalQuery += " AND ciudad.id=?"
			args = append(args, filters["pharmacy"])
		}
	} else if filters["pharmacy"] != 0 {
		finalQuery = " WHERE ciudad.id=?"
		args = []interface{}{filters["pharmacy"]}
	}
	log.Println(len(args))
	args = append(args, filters["limit"])
	args = append(args, filters["offset"])
	if len(args) == 0 {
		args = []interface{}{filters["limit"], filters["offset"]}
	}
	finalQuery += " LIMIT ? OFFSET ?"
	return finalQuery, &args
}

func GetDrugNames(db *sql.DB, ctx context.Context, id int32) (*[]string, error) {
	query := `SELECT monodrogas.nombre FROM monodrogas 
            JOIN medic_monodrogas
            ON medic_monodrogas.monodroga_id= monodrogas.id 
            JOIN medicamentos ON medicamentos.id = medic_monodrogas.medicamento_id
            WHERE medicamentos.id=?`

	var names []string
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		log.Println("Error")
		return &names, err
	}
	defer rows.Close()

	for rows.Next() {
		name := ""
		err := rows.Scan(&name)
		if err != nil {
			log.Println("Error")
			return &names, err
		}
		log.Printf("storing item: %+v", name)
		names = append(names, name)
	}
	return &names, nil
}

func handleMedicinePagination(db *sql.DB, ctx context.Context, filters map[string]int) bool {
	nextOffset := filters["limit"] + filters["offset"]
	sql, err := os.ReadFile(env.GetString("MED_COUNT", "../..internal/store/querys/medicines_count.sql"))
	if err != nil {
		log.Println(err)
		return false
	}

	query := string(sql)
	var args []interface{}

	if filters["drug"] != 0 {
		query += `WHERE monodrogas.id=?`
		args = []interface{}{filters["drug"]}
		if filters["pharmacy"] != 0 {
			query += " AND ciudad.id=?"
			args = append(args, filters["pharmacy"])
		}
	} else if filters["pharmacy"] != 0 {
		query = " WHERE ciudad.id=?"
		args = []interface{}{filters["pharmacy"]}
	}

	hasNextPage := false
	row := db.QueryRowContext(ctx, query, args...)
	var count int
	err = row.Scan(&count)
	if err != nil {
		log.Println("Error getting count: ", err)
	}

	if nextOffset < count {
		hasNextPage = true
	}
	return hasNextPage

}
