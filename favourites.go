package bete

import (
	"github.com/pkg/errors"
)

type FavouriteRepository interface {
	// Find searches for a favourite given a user ID and the name of the favourite,
	// returning the saved query if it exists and an empty string otherwise.
	Find(userID int, name string) string
	Put(userID int, name, query string) error
	List(userID int) ([]string, error)
}

type SQLFavouriteRepository struct {
	DB Queryable
}

func (r SQLFavouriteRepository) Find(userID int, text string) string {
	var query string
	err := r.DB.QueryRow("select query from favourites where user_id = $1 and name = $2", userID, text).Scan(&query)
	if err != nil {
		return ""
	}
	return query
}

func (r SQLFavouriteRepository) Put(userID int, name, query string) error {
	_, err := r.DB.Exec(`insert into favourites (user_id, name, query)
values ($1, $2, $3)
on conflict (user_id, name) do update set query = excluded.query`, userID, name, query)
	if err != nil {
		return errors.Wrap(err, "error putting new favourite")
	}
	return nil
}

func (r SQLFavouriteRepository) List(userID int) ([]string, error) {
	rows, err := r.DB.Query("select name from favourites where user_id = $1", userID)
	if err != nil {
		return nil, errors.Wrap(err, "error querying favourites")
	}
	defer rows.Close()
	var favourites []string
	for rows.Next() {
		var f string
		if err := rows.Scan(&f); err != nil {
			return favourites, errors.Wrap(err, "error scanning row")
		}
		favourites = append(favourites, f)
	}
	if err := rows.Err(); err != nil {
		return favourites, errors.Wrap(err, "error iterating rows")
	}
	return favourites, nil
}
