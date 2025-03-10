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
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Lastname    string `json:"lastname"`
	Birthday    string `json:"birthday"`
	C_ID        string `json:"c_id"`
	Role        string `json:"role"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phonenumber"`
}

type EmployeesStore struct {
	db *sql.DB
}

func (s *EmployeesStore) RegisterEmployee(ctx context.Context, e *Employee, u *User, branchID int64, date string) error {
	employeeQuery := `INSERT INTO empleados(nombre,apellido,fecha_nacimiento,
            cedula,cargo,direccion,telefono)
          VALUES(?,?,?,?,?,?,?) RETURNING id`
	var employeeID int
	err := s.db.QueryRowContext(ctx, employeeQuery, e.Name, e.Lastname, e.Birthday,
		e.C_ID, e.Role, e.Address, e.PhoneNumber).Scan(&employeeID)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Query completed for new row with id: %v", employeeID)

	userQuery := `INSERT INTO usuarios(correo,contrasena,codempleado)
          VALUES(?,?,?) RETURNING id`
	var userID int
	err = s.db.QueryRowContext(ctx, userQuery, u.Email, u.Password, employeeID).Scan(&userID)

	if err != nil {
		query := `DELETE FROM empleados WHERE id ?`
		_, nErr := s.db.ExecContext(ctx, query, employeeID)
		if nErr != nil {
			log.Println(err)
		}
		return err
	}
	log.Printf("Query completed for new user with id: %v", userID)
	query := `INSERT INTO rotacion(empleado_id,sucursal_id,fecha_inicio,observaciones)
              VALUES(?,?,?,?)`
	err = s.db.QueryRowContext(ctx, query, employeeID, branchID, date, "ingreso").Scan(&userID)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (s *EmployeesStore) Create(ctx context.Context, e *Employee) error {
	query := `INSERT INTO empleados(nombre,apellido,fecha_nacimiento,
            cedula,cargo,direccion,telefono)
          VALUES(?,?,?,?,?,?,?) RETURNING id`

	var id int
	err := s.db.QueryRowContext(ctx, query, e.Name, e.Lastname, e.Birthday,
		e.C_ID, e.Role, e.Address, e.PhoneNumber).Scan(&id)
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
	err := s.db.QueryRowContext(ctx, query, id).Scan(&e.ID, &e.Name, &e.Lastname, &e.Role, &e.Birthday, &e.Address, &e.PhoneNumber, &e.C_ID)
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
		err := rows.Scan(&item.ID, &item.Name, &item.Lastname, &item.Role, &item.Birthday, &item.Address, &item.PhoneNumber, &item.C_ID)
		if err != nil {
			log.Println("Error")
			return &employees, err
		}
		log.Printf("storing item: %+v", item)
		employees = append(employees, item)
	}
	return &employees, nil
}

func (s *EmployeesStore) GetFiltered(ctx context.Context, limit int, offset int, branch string) (*[]EmployeeView, bool, error) {
	sql, err := os.ReadFile(env.GetString("EMP_Q", "../../internal/store/querys/employees_view.sql"))
	if err != nil {
		log.Println(err)
		return nil, false, err
	}
	params, args := handleEmployeeFilters(branch, limit, offset)
	query := string(sql) + params
	log.Println(query)
	var employees []EmployeeView
	rows, err := s.db.QueryContext(ctx, query, *args...)
	if err != nil {
		log.Println(err)
		return nil, false, err
	}
	defer rows.Close()

	for rows.Next() {
		item := EmployeeView{}
		err := rows.Scan(&item.ID, &item.Name, &item.Lastname, &item.C_ID, &item.Address, &item.PhoneNumber,
			&item.Email, &item.Birthday, &item.Role)
		if err != nil {
			log.Println(err)
			return &employees, false, err
		}
		log.Printf("storing item: %+v", item)
		employees = append(employees, item)

	}
	hasNextPage := handleEmpPagination(s.db, ctx, limit+offset, branch)
	return &employees, hasNextPage, err
}

func handleEmployeeFilters(branch string, limit int, offset int) (string, *[]interface{}) {
	finalQuery := ""

	if branch == "" {
		finalQuery = `LIMIT ? OFFSET ?`
		return finalQuery, &[]interface{}{limit, offset}
	}
	finalQuery = `JOIN rotacion ON rotacion.empleado_id = empleados.id
                JOIN farmacia_sucursal ON farmacia_sucursal.id=rotacion.sucursal_id
                JOIN ciudad ON ciudad.id=farmacia_sucursal.ciudad_id
                WHERE ciudad.nombre= ? LIMIT ? OFFSET ?`

	return finalQuery, &[]interface{}{branch, limit, offset}
}

func handleEmpPagination(db *sql.DB, ctx context.Context, nextOffset int, branch string) bool {
	var query string
	if branch != "" {
		query = `SELECT COUNT(*) FROM empleados JOIN usuarios 
    ON usuarios.codempleado = empleados.id 
    JOIN rotacion ON rotacion.empleado_id = empleados.id
    JOIN farmacia_sucursal ON farmacia_sucursal.id=rotacion.sucursal_id
    JOIN ciudad ON ciudad.id=farmacia_sucursal.ciudad_id
    WHERE ciudad.nombre= ?`
	} else {
		query = `SELECT COUNT(*) FROM 
    empleados JOIN usuarios ON usuarios.codempleado = empleados.id`
	}

	row := db.QueryRowContext(ctx, query, branch)

	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Println("Error getting count: ", err)
		return false
	}

	if nextOffset < count {
		return true
	}
	return false
}
