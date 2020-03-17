package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/yi-jiayu/datamall/v3"
)

var db *sql.DB

func createStops() error {
	log.Println("creating table stops")
	_, err := db.Exec(`create table if not exists stops (id text primary key, road text, description text, location point)`)
	if err != nil {
		return fmt.Errorf("error creating table stops: %w", err)
	}
	return nil
}

func syncStops(accountKey string) error {
	log.Println("syncing data in stops table")
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

func createFavourites() error {
	log.Println("creating table favourites")
	_, err := db.Exec(`create table if not exists favourites
(
    user_id integer not null,
    name    text    not null,
    query   text    not null
)`)
	if err != nil {
		return fmt.Errorf("error creating table favourites: %w", err)
	}
	return nil
}

var sync = flag.Bool("sync", false, "whether to sync data from datamall")

func main() {
	flag.Parse()

	var err error
	databaseURL := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
	}

	if err := createStops(); err != nil {
		log.Fatal(err)
	}
	if err := createFavourites(); err != nil {
		log.Fatal(err)
	}

	if *sync {
		accountKey := os.Getenv("DATAMALL_ACCOUNT_KEY")
		if accountKey == "" {
			log.Fatal("DATAMALL_ACCOUNT_KEY environment variable not set")
		}

		if err := syncStops(accountKey); err != nil {
			log.Fatal(err)
		}
	}
}
