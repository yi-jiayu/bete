package bete

import (
	"database/sql"
)

type FavouriteRepository interface {
	// FindByUserAndText searches for a favourite given a user ID and the name of the favourite,
	// returning the saved query if it exists and an empty string otherwise.
	FindByUserAndText(userID int, name string) string
}

type PostgresFavouriteRepository struct {
	DB *sql.DB
}

func (r PostgresFavouriteRepository) FindByUserAndText(userID int, text string) string {
	var query string
	err := r.DB.QueryRow("select query from favourites where user_id = $1 and name = $2", userID, text).Scan(&query)
	if err != nil {
		return ""
	}
	return query
}
