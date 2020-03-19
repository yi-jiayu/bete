package bete

import (
	"database/sql"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Queryable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
