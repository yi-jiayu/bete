package bete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLFavouriteRepository_Find(t *testing.T) {
	tx := getDatabaseTx()
	defer tx.Rollback()

	userID := randomInt64ID()
	name := "SUTD"
	query := "96049 5 24"
	_, err := tx.Exec(`insert into favourites (user_id, name, query) values ($1, $2, $3)`, userID, name, query)
	if err != nil {
		t.Fatal(err)
	}

	repo := SQLFavouriteRepository{DB: tx}
	actual := repo.Find(userID, name)
	assert.Equal(t, query, actual)
}

func TestSQLFavouriteRepository_Put(t *testing.T) {
	tx := getDatabaseTx()
	defer tx.Rollback()

	userID := randomInt64ID()
	name := "SUTD"
	query := "96049 5 24"

	repo := SQLFavouriteRepository{DB: tx}
	err := repo.Put(userID, name, query)
	assert.NoError(t, err)
	var actual string
	err = tx.QueryRow("select query from favourites where user_id = $1 and name = $2", userID, name).Scan(&actual)
	assert.Equal(t, query, actual)
}

func TestSQLFavouriteRepository_List(t *testing.T) {
	tx := getDatabaseTx()
	defer tx.Rollback()

	userID := randomInt64ID()
	names := []string{"SUTD", "Paya Lebar MRT"}
	_, err := tx.Exec(
		`insert into favourites (user_id, name, query)
values ($1, $2, $3),
       ($4, $5, $6),
       ($7, $8, $9)`,
		userID, names[0], "96049 5 24",
		userID, names[1], "81111",
		456, "UIC Bldg", "03129",
	)
	if err != nil {
		t.Fatal(err)
	}

	repo := SQLFavouriteRepository{DB: tx}
	actual, err := repo.List(userID)
	assert.NoError(t, err)
	assert.Equal(t, names, actual)
}

func TestSQLFavouriteRepository_Delete(t *testing.T) {
	tx := getDatabaseTx()
	defer tx.Rollback()

	userID := randomInt64ID()
	name := "SUTD"
	query := "96049 5 24"
	_, err := tx.Exec(`insert into favourites (user_id, name, query) values ($1, $2, $3)`, userID, name, query)
	if err != nil {
		t.Fatal(err)
	}

	repo := SQLFavouriteRepository{DB: tx}
	err = repo.Delete(userID, name)
	assert.NoError(t, err)

	favourites, err := repo.List(userID)
	assert.NoError(t, err)
	assert.Len(t, favourites, 0)
}
