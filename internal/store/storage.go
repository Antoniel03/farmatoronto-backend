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
		GetFiltered(context.Context, int, int, string, string) (*[]MedicineView, bool, error)
	}

	Employees interface {
		Create(context.Context, *Employee) error
		GetByID(context.Context, string) (*Employee, error)
		GetAll(context.Context) (*[]Employee, error)
		GetFiltered(context.Context, int, int, string) (*[]EmployeeView, bool, error)
		RegisterEmployee(context.Context, *Employee, *User, int64, string) error
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
		GetPaginated(context.Context, int, int) (*[]Lab, bool, error)
	}

	Branches interface {
		Create(context.Context, *Branch) error
		GetByID(context.Context, string) (*Branch, error)
		GetAll(context.Context) (*[]Branch, error)
		GetPaginated(context.Context, int, int) (*[]Branch, bool, error)
	}

	Drugs interface {
		Create(context.Context, *Drug) error
		GetAll(context.Context) (*[]Drug, error)
		GetPaginated(context.Context, int, int) (*[]Drug, bool, error)
	}
	Actions interface {
		Create(context.Context, string) error
		GetAll(context.Context) (*[]Action, error)
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
		Actions:   &ActionsStore{db},
	}
}
