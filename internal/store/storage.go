package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Medicines interface {
		Create(context.Context, *Medicine) error
		GetByID(context.Context, string) (*Medicine, error)
		GetAll(context.Context) (*[]Medicine, error)
	}

	Employees interface {
		Create(context.Context, *Employee) error
	}

	Users interface {
		Create(context.Context, *User) error
		GetByEmail(context.Context, string) (*User, error)
		GetByID(context.Context, string) (*User, error)
		GetAll(context.Context) (*[]User, error)
	}
}

func NewSQLiteStorage(db *sql.DB) Storage {
	return Storage{
		Medicines: &MedicinesStore{db},
		Employees: &EmployeesStore{db},
		Users:     &UsersStore{db},
	}
}
