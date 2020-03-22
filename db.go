package bete

import (
	"database/sql"
)

// Conn is an interface containing methods available on sql.DB and sql.Tx objects.
type Conn interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
