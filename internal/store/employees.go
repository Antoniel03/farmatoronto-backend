package store

import (
	"context"
	"database/sql"
)

type Employee struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Lastname    string `json:"lastname"`
	Birthday    string `json:"birthday`
	Direction   string `json:"direction"`
	PhoneNumber string `json:"phonenumber"`
	Email       string `json:"email"`
}

type EmployeesStore struct {
	db *sql.DB
}

func (s *EmployeesStore) Create(ctx context.Context, e *Employee) error {
	return nil
}
