package interfaces

import (
	"errors"
)

var ErrNoRows = errors.New("No rows")

type DbHandler interface {
	Execute(statement string, params ...interface{}) (Result, error)
	Query(statement string, params ...interface{}) (Rows, error)
	QueryRow(statement string, params ...interface{}) Row
}

type Row interface {
	Scan(dest ...interface{}) error
}

type Rows interface {
	Scan(dest ...interface{}) error
	Next() bool
}

type Result interface {
	LastInsertId() (int64, error)
}

type DbRepo struct {
	dbHandlers map[string]DbHandler
	dbHandler  DbHandler
}
