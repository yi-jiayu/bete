package bete

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type BusStop struct {
	ID          string
	Description string
	RoadName    string
}

type BusStopRepository interface {
	Find(id string) (BusStop, error)
}

type PostgresBusStopRepository struct {
	db *sql.DB
}

func NewPostgresBusStopRepository(url string) (PostgresBusStopRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return PostgresBusStopRepository{}, errors.Wrap(err, "error opening postgres database")
	}
	return PostgresBusStopRepository{db: db}, nil
}

func (r PostgresBusStopRepository) Find(id string) (BusStop, error) {
	var stop BusStop
	err := r.db.QueryRow("select id, description, road from stops where id = $1", id).Scan(&stop.ID, &stop.Description, &stop.RoadName)
	if err == sql.ErrNoRows {
		return stop, Error("not found")
	} else if err != nil {
		return stop, errors.Wrap(err, "error querying bus stop")
	}
	return stop, nil
}
