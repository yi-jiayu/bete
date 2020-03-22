package bete

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

type BusStop struct {
	ID          string
	Description string
	RoadName    string
	Location    Location
}

// NearbyBusStop represents how far away a bus stop is.
type NearbyBusStop struct {
	BusStop
	Distance float32
}

// Location represents a latitude and longitude coordinate.
type Location struct {
	Latitude  float32
	Longitude float32
}

func (l *Location) Scan(src interface{}) error {
	point := src.(string)
	n, err := fmt.Sscanf(point, "(%f,%f)", &l.Longitude, &l.Latitude)
	if err != nil {
		return errors.Wrap(err, "error scanning point")
	}
	if n != 2 {
		return errors.New("not enough values")
	}
	return nil
}

type BusStopRepository interface {
	Find(id string) (BusStop, error)

	// Nearby returns up to limit bus stops within radius km of the point specified by lat and lon.
	Nearby(lat, lon, radius float32, limit int) ([]NearbyBusStop, error)

	// Search searches for stops with a text query.
	Search(query string, limit int) ([]BusStop, error)
}

type SQLBusStopRepository struct {
	DB Conn
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

func (r SQLBusStopRepository) Nearby(lat, lon, radius float32, limit int) ([]NearbyBusStop, error) {
	location := fmt.Sprintf("(%f, %f)", lon, lat)
	rows, err := r.DB.Query(
		`select id, road, description, location::text, (location <@> $1) * 1.609344 as distance
from stops
where (location <@> $1) * 1.609344 < $2
order by distance
limit $3`,
		location,
		radius,
		limit,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error querying nearby bus stops")
	}
	defer rows.Close()
	var nearby []NearbyBusStop
	for rows.Next() {
		var n NearbyBusStop
		if err := rows.Scan(&n.ID, &n.RoadName, &n.Description, &n.Location, &n.Distance); err != nil {
			return nearby, errors.Wrap(err, "error scanning row")
		}
		nearby = append(nearby, n)
	}
	if err := rows.Err(); err != nil {
		return nearby, errors.Wrap(err, "error iterating rows")
	}
	return nearby, nil
}

func (r SQLBusStopRepository) Search(query string, limit int) ([]BusStop, error) {
	var rows *sql.Rows
	var err error
	if query != "" {
		rows, err = r.DB.Query(
			`select id, road, description, location::text
from stops
where tokens @@ to_tsquery($1)
order by ts_rank(tokens, to_tsquery($1)) desc
limit $2;`,
			query,
			limit,
		)
	} else {
		rows, err = r.DB.Query(
			`select id, road, description, location::text
from stops
order by id
limit $1;`,
			limit,
		)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "error searching bus stops (query: %q)", query)
	}
	defer rows.Close()
	var matches []BusStop
	for rows.Next() {
		var s BusStop
		if err := rows.Scan(&s.ID, &s.RoadName, &s.Description, &s.Location); err != nil {
			return matches, errors.Wrap(err, "error scanning row")
		}
		matches = append(matches, s)
	}
	if err := rows.Err(); err != nil {
		return matches, errors.Wrap(err, "error iterating rows")
	}
	return matches, nil
}
