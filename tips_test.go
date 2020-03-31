package bete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLTipsRepository_Get(t *testing.T) {
	tx := getDatabaseTx()
	defer tx.Rollback()

	tips := []interface{}{"Tip A", "Tip B", "Tip C"}
	_, err := tx.Exec(`insert into tips (content) values ($1), ($2), ($3)`, tips...)
	if err != nil {
		t.Fatal(err)
	}

	repo := SQLTipsRepository{DB: tx}
	var tip string
	tip, err = repo.Get(0)
	assert.NoError(t, err)
	assert.Equal(t, "Tip A", tip)
	tip, err = repo.Get(5)
	assert.NoError(t, err)
	assert.Equal(t, "Tip C", tip)
	tip, err = repo.Get(1234)
	assert.NoError(t, err)
	assert.Equal(t, "Tip B", tip)
}
