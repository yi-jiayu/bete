package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/yi-jiayu/datamall/v3"
)

func syncStops(db *sql.DB, dm datamall.APIClient) (int, error) {
	txn, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error beginning transaction: %w", err)
	}
	stmt, err := txn.Prepare(`insert into stops
values ($1, $2, $3, $4)
on conflict (id) do update set road        = excluded.road,
                               description = excluded.description,
                               location    = excluded.location`)
	if err != nil {
		return 0, fmt.Errorf("error preparing statement: %w", err)
	}
	offset := 0
	for {
		stops, err := dm.GetBusStops(offset)
		if err != nil {
			return offset, fmt.Errorf("error getting bus stops: %w", err)
		}
		count := len(stops.Value)
		if count == 0 {
			break
		}
		for _, stop := range stops.Value {
			location := fmt.Sprintf("(%f, %f)", stop.Longitude, stop.Latitude)
			_, err := stmt.Exec(stop.BusStopCode, stop.RoadName, stop.Description, location)
			if err != nil {
				return offset, fmt.Errorf("error inserting stop: %w", err)
			}
		}
		offset += count
		log.Printf("inserted %d stops", count)
	}
	err = stmt.Close()
	if err != nil {
		return offset, fmt.Errorf("error closing statement: %w", err)
	}
	err = txn.Commit()
	if err != nil {
		return offset, fmt.Errorf("error committing txn: %w", err)
	}
	return offset, nil
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		panic(err)
	}

	accountKey := os.Getenv("DATAMALL_ACCOUNT_KEY")
	if accountKey == "" {
		panic("DATAMALL_ACCOUNT_KEY environment variable not set")
	}
	dm := datamall.NewDefaultClient(accountKey)

	log.Println("syncing bus stop data from datamall")
	n, err := syncStops(db, dm)
	if err != nil {
		panic(err)
	}
	log.Printf("synced %d bus stops from datamall", n)
}
