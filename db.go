package bete

import (
	"database/sql"
)

const ErrNotFound = Error("not found")

type Error string

func (e Error) Error() string {
	return string(e)
}

type Queryable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
