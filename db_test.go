package bete

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func getDatabaseTx() *sql.Tx {
	if db == nil {
		var err error
		db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			panic(err)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	return tx
}
