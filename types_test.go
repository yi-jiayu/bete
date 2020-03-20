package bete

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		query Query
		err   error
	}{
		{
			name: "bus stop code only",
			text: "96049",
			query: Query{
				Stop:   "96049",
				Filter: []string{},
			},
		},
		{
			name: "bus stop code and filter",
			text: "96049 5 24",
			query: Query{
				Stop:   "96049",
				Filter: []string{"5", "24"},
			},
		},
		{
			name:  "does not start with a bus stop code",
			text:  "ABCDE 5 24",
			query: Query{},
			err:   ErrQueryDoesNotStartWithBusStopCode,
		},
		{
			name:  "invalid characters",
			text:  "12345 5! '2'",
			query: Query{},
			err:   ErrQueryContainsInvalidCharacters,
		},
		{
			name:  "too long",
			text:  "12345 24 28 43 70 76 134 135 137 154 155",
			query: Query{},
			err:   ErrQueryTooLong,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := ParseQuery(tt.text)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.query, query)
		})
	}
}

func TestQuery_Canonical(t *testing.T) {
	query1 := Query{Stop: "96049"}
	assert.Equal(t, "96049", query1.Canonical())
	query2 := Query{Stop: "96049", Filter: []string{"5", "24"}}
	assert.Equal(t, "96049 5 24", query2.Canonical())
}

func TestQuery_MarshalJSON(t *testing.T) {
	query1 := Query{Stop: "96049"}
	json1, err := json.Marshal(&query1)
	assert.NoError(t, err)
	assert.Equal(t, `"96049"`, string(json1))
	query2 := Query{Stop: "96049", Filter: []string{"5", "24"}}
	json2, err := json.Marshal(&query2)
	assert.NoError(t, err)
	assert.Equal(t, `"96049 5 24"`, string(json2))
}

func TestQuery_UnmarshalJSON(t *testing.T) {
	var query Query
	var err error
	err = json.Unmarshal([]byte(`"96049"`), &query)
	assert.NoError(t, err)
	assert.Equal(t, Query{Stop: "96049", Filter: []string{}}, query)
	err = json.Unmarshal([]byte(`"96049 5 24"`), &query)
	assert.NoError(t, err)
	assert.Equal(t, Query{Stop: "96049", Filter: []string{"5", "24"}}, query)
}
