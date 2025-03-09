package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Medicines interface {
		Create(context.Context, *Medicine, *MedicineExtraData) error
		GetByID(context.Context, string) (*Medicine, error)
		GetAll(context.Context) (*[]Medicine, error)
		GetPaginated(context.Context, int, int) (*[]Medicine, error)
		GetFiltered(context.Context, int, int, string, string) (*[]MedicineView, error)
	}

	Employees interface {
		Create(context.Context, *Employee) error
		GetByID(context.Context, string) (*Employee, error)
		GetAll(context.Context) (*[]Employee, error)
		GetFiltered(context.Context, int, int, string) (*[]EmployeeView, error)
	}

	Users interface {
		Create(context.Context, *User) error
		GetByEmail(context.Context, string) (*User, error)
		GetByID(context.Context, string) (*User, error)
		GetAll(context.Context) (*[]User, error)
		GetPaginated(context.Context, int, int) (*[]User, error)
		GetLoginData(context.Context, string) (*User, *Employee, error)
	}

	Labs interface {
		Create(context.Context, *Lab) error
		GetByID(context.Context, string) (*Lab, error)
		GetAll(context.Context) (*[]Lab, error)
		GetPaginated(context.Context, int, int) (*[]Lab, error)
	}

	Branches interface {
		Create(context.Context, *Branch) error
		GetByID(context.Context, string) (*Branch, error)
		GetAll(context.Context) (*[]Branch, error)
		GetPaginated(context.Context, int, int) (*[]Branch, error)
	}

	Drugs interface {
		Create(context.Context, *Drugs) error
		// GetByID(context.Context, string) (*Drugs, error)
		// GetAll(context.Context) (*[]Branch, error)
		// GetPaginated(context.Context, int, int) (*[]Drugs, error)
	}
}

func NewSQLiteStorage(db *sql.DB) Storage {
	return Storage{
		Medicines: &MedicinesStore{db},
		Employees: &EmployeesStore{db},
		Users:     &UsersStore{db},
		Labs:      &LabsStore{db},
		Branches:  &BranchesStore{db},
		Drugs:     &DrugsStore{db},
	}
}
