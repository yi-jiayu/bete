package bete

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type BusStop struct {
	ID          string
	Description string
	RoadName    string
}

type BusStopRepository interface {
	Find(id string) (BusStop, error)
}

type SQLBusStopRepository struct {
	DB Queryable
}

func (r SQLBusStopRepository) Find(id string) (BusStop, error) {
	var stop BusStop
	err := r.DB.QueryRow("select id, description, road from stops where id = $1", id).Scan(&stop.ID, &stop.Description, &stop.RoadName)
	if err == sql.ErrNoRows {
		return stop, ErrNotFound
	} else if err != nil {
		return stop, errors.Wrap(err, "error querying bus stop")
	}
	return stop, nil
}
