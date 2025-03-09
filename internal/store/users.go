package store

import (
	"context"
	"database/sql"
	"log"
)

type User struct {
	ID         int64  `json:"id"`
	EmployeeID int64  `json:"employee_id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, u *User) error {
	query := `INSERT INTO usuarios(correo,contrasena,codempleado)
          VALUES(?,?,?) RETURNING id`
	var id int
	err := s.db.QueryRowContext(ctx, query, u.Email, u.Password, u.EmployeeID).Scan(&id)
	if err != nil {
		return err
	}
	log.Printf("Query completed for new user with id: %v", id)
	return nil
}

func (s *UsersStore) GetAll(ctx context.Context) (*[]User, error) {
	query := `SELECT * FROM usuarios`
	var users []User

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		item := User{}
		err := rows.Scan(&item.ID, &item.Email, &item.Password, &item.EmployeeID)
		if err != nil {
			return &users, err
		}
		users = append(users, item)
	}
	return &users, nil
}

func (s *UsersStore) GetByID(ctx context.Context, id string) (*User, error) {
	query := `SELECT * FROM usuarios WHERE id=?`

	u := User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Email, &u.Password, &u.EmployeeID)
	if err != nil {
		return nil, err
	}
	log.Printf("Query completed for the requested item\n%+v", u)
	return &u, nil
}

func (s *UsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT * FROM usuarios WHERE correo=?`

	u := User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Password, &u.EmployeeID)
	if err != nil {
		return nil, err
	}
	log.Printf("Query completed for the requested item\n%+v", u)
	return &u, nil
}

func (s *UsersStore) GetPaginated(ctx context.Context, limit int, offset int) (*[]User, error) {
	query := `SELECT * FROM usuarios LIMIT ? OFFSET ?`
	var users []User

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		item := User{}
		err := rows.Scan(&item.ID, &item.Email, &item.Password, &item.EmployeeID)
		if err != nil {
			return &users, err
		}
		users = append(users, item)
	}
	return &users, nil
}

func (s *UsersStore) GetLoginData(ctx context.Context, email string) (*User, *Employee, error) {
	query := `SELECT usuarios.id,usuarios.correo,usuarios.contrasena, 
            empleados.cargo, empleados.nombre, empleados.apellido
            FROM usuarios JOIN empleados ON empleados.id=usuarios.codempleado
            WHERE usuarios.correo=?`

	e := Employee{}
	u := User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email,
		&u.Password, &e.Role, &e.Name, &e.Lastname)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	log.Printf("Query completed for the requested item\n%+v", u)
	return &u, &e, nil
}
