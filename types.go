package bete

import (
	"encoding/json"
	"regexp"
	"strings"
)

const (
	ErrNotFound                         = Error("not found")
	ErrQueryDoesNotStartWithBusStopCode = Error("query does not start with a bus stop code")
	ErrQueryContainsInvalidCharacters   = Error("query contains invalid characters")
	ErrQueryTooLong                     = Error("query too long")
)

const MaxQueryLength = 20

var (
	queryStartsWithBusStopCodeRegexp       = regexp.MustCompile(`^\d{5}(?:\s|$)`)
	queryContainsOnlyValidCharactersRegexp = regexp.MustCompile(`^[0-9A-Za-z\s]+$`)
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type CallbackData struct {
	Type   string   `json:"t"`
	StopID string   `json:"b,omitempty"`
	Filter []string `json:"s,omitempty"`
	Format string   `json:"f,omitempty"`
	Name   string   `json:"n,omitempty"`
	Query  *Query   `json:"q,omitempty"`
}

func (c CallbackData) Encode() string {
	JSON, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(JSON)
}

type Query struct {
	Stop   string
	Filter []string
}

func (q *Query) UnmarshalJSON(data []byte) error {
	var text string
	err := json.Unmarshal(data, &text)
	if err != nil {
		return err
	}
	parsed, err := ParseQuery(text)
	if err != nil {
		return err
	}
	q.Stop = parsed.Stop
	q.Filter = parsed.Filter
	return nil
}

func (q *Query) MarshalJSON() ([]byte, error) {
	return json.Marshal(q.Canonical())
}

func ParseQuery(text string) (Query, error) {
	if ok := queryStartsWithBusStopCodeRegexp.MatchString(text); !ok {
		return Query{}, ErrQueryDoesNotStartWithBusStopCode
	}
	if ok := queryContainsOnlyValidCharactersRegexp.MatchString(text); !ok {
		return Query{}, ErrQueryContainsInvalidCharacters
	}
	if len(text) > MaxQueryLength {
		return Query{}, ErrQueryTooLong
	}
	parts := strings.Fields(text)
	return Query{
		Stop:   parts[0],
		Filter: parts[1:],
	}, nil
}

func (q *Query) Canonical() string {
	return strings.Join(append([]string{q.Stop}, q.Filter...), " ")
}
