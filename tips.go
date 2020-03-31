package bete

import (
	"github.com/pkg/errors"
)

type TipsRepository interface {
	Get(seed int) (string, error)
}

type SQLTipsRepository struct {
	DB Conn
}

func (s SQLTipsRepository) Get(seed int) (string, error) {
	var tip string
	err := s.DB.QueryRow(`select content
from tips
order by id
limit 1
offset
mod($1, (select count(*) from tips))`, seed).Scan(&tip)
	if err != nil {
		return "", errors.Wrapf(err, "error querying tip")
	}
	return tip, nil
}
