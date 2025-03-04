package store

import (
	"context"
	"database/sql"
	"log"
)

type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, u *User) error {
	query := `INSERT INTO usuarios(correo,contrasena,tipo_usuario)
          VALUES(?,?,?) RETURNING id`
	var id int
	err := s.db.QueryRowContext(ctx, query, u.Email, u.Password, u.UserType).Scan(&id)
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
		err := rows.Scan(&item.Id, &item.Email, &item.Password, &item.UserType)
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
	err := s.db.QueryRowContext(ctx, query, id).Scan(&u.Id, &u.Email, &u.Password, &u.UserType)
	if err != nil {
		return nil, err
	}
	log.Printf("Query completed for the requested item\n%+v", u)
	return &u, nil
}

func (s *UsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT * FROM usuarios WHERE correo=?`

	u := User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(&u.Id, &u.Email, &u.Password, &u.UserType)
	if err != nil {
		return nil, err
	}
	log.Printf("Query completed for the requested item\n%+v", u)
	return &u, nil
}
