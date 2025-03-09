package store

import (
	"context"
	"database/sql"
	"github.com/Antoniel03/farmatoronto-backend/internal/env"
	"log"
	"os"
)

type EmployeeView struct {
	Employee
	Email string `json:"email"`
}

type Employee struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Lastname    string `json:"lastname"`
	Birthday    string `json:"birthday"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phonenumber"`
}

type EmployeesStore struct {
	db *sql.DB
}

func (s *EmployeesStore) Create(ctx context.Context, e *Employee) error {
	query := `INSERT INTO empleados(nombre,apellido,fecha_nacimiento,
            direccion,telefono)
          VALUES(?,?,?,?,?) RETURNING id`

	var id int
	err := s.db.QueryRowContext(ctx, query, e.Name, e.Lastname, e.Birthday,
		e.Address, e.PhoneNumber).Scan(&id)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Query completed for new medicine with id: %v", id)
	return nil
}

func (s *EmployeesStore) GetByID(ctx context.Context, id string) (*Employee, error) {
	query := `SELECT * FROM empleados WHERE id=?`

	e := Employee{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&e.ID, &e.Name, &e.Lastname, &e.Birthday, &e.Address, &e.PhoneNumber)
	if err != nil {
		return nil, err
	}
	log.Printf("Query completed for the requested item\n%+v", e)
	return &e, nil
}

func (s *EmployeesStore) GetAll(ctx context.Context) (*[]Employee, error) {
	query := `SELECT * FROM empleados`
	var employees []Employee
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := Employee{}
		err := rows.Scan(&item.ID, &item.Name, &item.Lastname, &item.Birthday, &item.Address, &item.PhoneNumber)
		if err != nil {
			log.Println("Error")
			return &employees, err
		}
		log.Printf("storing item: %+v", item)
		employees = append(employees, item)
	}
	return &employees, nil
}

func (s *EmployeesStore) GetFiltered(ctx context.Context, limit int, offset int, branch string) (*[]EmployeeView, error) {
	sql, err := os.ReadFile(env.GetString("EMP_Q", "../..internal/store/querys/employees_view.sql"))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	params, args := handleEmployeeFilters(branch, limit, offset)

	query := string(sql) + params
	log.Println(query)
	var employees []EmployeeView
	rows, err := s.db.QueryContext(ctx, query, *args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := EmployeeView{}
		err := rows.Scan(&item.ID, &item.Name, &item.Lastname, &item.Birthday,
			&item.Address, &item.PhoneNumber, &item.Email)
		if err != nil {
			log.Println(err)
			return &employees, err
		}
		log.Printf("storing item: %+v", item)
		employees = append(employees, item)
	}
	return &employees, nil
}

func handleEmployeeFilters(branch string, limit int, offset int) (string, *[]interface{}) {
	finalQuery := ""

	if branch == "" {
		return finalQuery, &[]interface{}{limit, offset}
	}
	finalQuery = `JOIN rotacion on rotacion.empleado_id = empleados.id
                JOIN farmacia_sucursal ON farmacia_sucursal.id=rotacion.sucursal_id
                  WHERE farmacia_sucursal.nombre= ? LIMIT ? OFFSET ?`

	return finalQuery, &[]interface{}{branch, limit, offset}
}
