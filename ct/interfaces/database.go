package interfaces

import (
	"errors"
)

const dbDateFormat = "2006-01-02 15:04:05"

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
	Close() error
}

type Result interface {
	LastInsertId() (int64, error)
}

type DbRepo struct {
	dbHandlers map[string]DbHandler
	dbHandler  DbHandler
}
