package repository

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Repositories struct {
	UserRepository *UserRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepository: NewRepository(db),
	}
}
