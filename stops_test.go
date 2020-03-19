package bete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresBusStopRepository_Find(t *testing.T) {
	tx := getDatabaseTx()
	defer tx.Rollback()

	stop := BusStop{
		ID:          "99999",
		Description: "Test Description",
		RoadName:    "Test Road",
	}
	_, err := tx.Exec(`insert into stops (id, road, description, location)
values ($1, $2, $3, '(0.123, 0.456)')`, stop.ID, stop.RoadName, stop.Description)
	if err != nil {
		t.Fatal(err)
	}

	repo := SQLBusStopRepository{DB: tx}
	actual, err := repo.Find(stop.ID)
	assert.NoError(t, err)
	assert.Equal(t, stop, actual)
}
