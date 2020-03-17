package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/yi-jiayu/datamall/v3"
)

func syncBusStops(accountKey, databaseURL string) error {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("error connecting to postgres: %w", err)
	}
	_, err = db.Exec(`create table if not exists stops (id text primary key, road text, description text, location point)`)
	if err != nil {
		return fmt.Errorf("error creating table stops: %w", err)
	}
	txn, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	stmt, err := txn.Prepare(`insert into stops
values ($1, $2, $3, $4)
on conflict (id) do update set road        = excluded.road,
                               description = excluded.description,
                               location    = excluded.location`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	dm := datamall.NewDefaultClient(accountKey)
	offset := 0
	for {
		stops, err := dm.GetBusStops(offset)
		if err != nil {
			return fmt.Errorf("error getting bus stops: %w", err)
		}
		count := len(stops.Value)
		if count == 0 {
			break
		}
		for _, stop := range stops.Value {
			location := fmt.Sprintf("(%f, %f)", stop.Latitude, stop.Longitude)
			_, err := stmt.Exec(stop.BusStopCode, stop.RoadName, stop.Description, location)
			if err != nil {
				return fmt.Errorf("error inserting stop: %w", err)
			}
		}
		offset += count
		log.Printf("inserted %d stops", count)
	}
	err = stmt.Close()
	if err != nil {
		return fmt.Errorf("error closing statement: %w", err)
	}
	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("error committing txn: %w", err)
	}
	return nil
}

func main() {
	accountKey := os.Getenv("DATAMALL_ACCOUNT_KEY")
	databaseURL := os.Getenv("DATABASE_URL")

	var err error
	err = syncBusStops(accountKey, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
}
