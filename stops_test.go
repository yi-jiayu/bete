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

func TestSQLBusStopRepository_Nearby(t *testing.T) {
	tx := getDatabaseTx()
	defer tx.Rollback()

	stmt, err := tx.Prepare(`insert into stops (id, road, description, location) values ($1, $2, $3, $4)`)
	if err != nil {
		t.Fatal(err)
	}
	stops := [][]interface{}{
		{"01319", "Kallang Rd", "Lavender Stn Exit A/ICA", "(103.863256,1.307574)"},
		{"01339", "Crawford St", "Bef Crawford Bridge", "(103.864263,1.307746)"}, // 0.11356564947243729 km away
		{"07371", "Lavender St", "Aft Kallang Rd", "(103.863501,1.309508)"},      // 0.21676780485189698 km away
	}
	for _, stop := range stops {
		_, err := stmt.Exec(stop...)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name     string
		radius   float32
		limit    int
		expected []NearbyBusStop
	}{
		{
			name:   "stops within 500 m",
			radius: 0.5,
			expected: []NearbyBusStop{
				{
					BusStop: BusStop{
						ID:          "01319",
						Description: "Lavender Stn Exit A/ICA",
						RoadName:    "Kallang Rd",
						Location:    Location{Latitude: 1.307574, Longitude: 103.863256}},
					Distance: 0,
				},
				{
					BusStop: BusStop{
						ID:          "01339",
						Description: "Bef Crawford Bridge",
						RoadName:    "Crawford St",
						Location:    Location{Latitude: 1.307746, Longitude: 103.864263},
					},
					Distance: 0.11356564947243729,
				},
				{
					BusStop: BusStop{
						ID:          "07371",
						Description: "Aft Kallang Rd",
						RoadName:    "Lavender St",
						Location:    Location{Latitude: 1.309508, Longitude: 103.863501},
					},
					Distance: 0.21676780485189698,
				},
			},
			limit: 50,
		},
		{
			name:   "up to 2 stops within 500 m",
			radius: 0.5,
			expected: []NearbyBusStop{
				{
					BusStop: BusStop{
						ID:          "01319",
						Description: "Lavender Stn Exit A/ICA",
						RoadName:    "Kallang Rd",
						Location:    Location{1.307574, 103.863256}},
					Distance: 0,
				},
				{
					BusStop: BusStop{
						ID:          "01339",
						Description: "Bef Crawford Bridge",
						RoadName:    "Crawford St",
						Location:    Location{1.307746, 103.864263},
					},
					Distance: 0.11356564947243729,
				},
			},
			limit: 2,
		},
	}

	var lat, lon float32 = 1.307574, 103.863256
	repo := SQLBusStopRepository{DB: tx}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := repo.Nearby(lat, lon, tt.radius, tt.limit)
			assert.NoError(t, err)
			assert.Len(t, actual, len(tt.expected))
			for i, nearby := range actual {
				assert.Equal(t, tt.expected[i].BusStop, nearby.BusStop)
				assert.InDelta(t, tt.expected[i].Distance, nearby.Distance, 0.001)
			}
		})
	}
}

func TestSQLBusStopRepository_Search(t *testing.T) {
	tx := getDatabaseTx()
	defer tx.Rollback()

	stmt, err := tx.Prepare(`insert into stops (id, road, description, location) values ($1, $2, $3, $4)`)
	if err != nil {
		t.Fatal(err)
	}
	stops := [][]interface{}{
		{"01319", "Kallang Rd", "Lavender Stn Exit A/ICA", "(103.863256,1.307574)"},
		{"01339", "Crawford St", "Bef Crawford Bridge", "(103.864263,1.307746)"},
		{"07371", "Lavender St", "Aft Kallang Rd", "(103.863501,1.309508)"},
	}
	for _, stop := range stops {
		_, err := stmt.Exec(stop...)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name     string
		query    string
		limit    int
		expected []BusStop
	}{
		{
			name:  "returns all matches",
			query: "lavender",
			limit: 50,
			expected: []BusStop{
				{
					ID:          "01319",
					Description: "Lavender Stn Exit A/ICA",
					RoadName:    "Kallang Rd",
					Location:    Location{Latitude: 1.307574, Longitude: 103.86326},
				},
				{
					ID:          "07371",
					Description: "Aft Kallang Rd",
					RoadName:    "Lavender St",
					Location:    Location{Latitude: 1.309508, Longitude: 103.8635},
				},
			},
		},
		{
			name:  "query with spaces",
			query: "lavender stn",
			limit: 50,
			expected: []BusStop{
				{
					ID:          "01319",
					Description: "Lavender Stn Exit A/ICA",
					RoadName:    "Kallang Rd",
					Location:    Location{Latitude: 1.307574, Longitude: 103.86326},
				},
			},
		},
		{
			name:  "returns all matches with limit",
			query: "lavender",
			limit: 1,
			expected: []BusStop{
				{
					ID:          "01319",
					Description: "Lavender Stn Exit A/ICA",
					RoadName:    "Kallang Rd",
					Location:    Location{Latitude: 1.307574, Longitude: 103.86326},
				},
			},
		},
		{
			name:  "returns all bus stops when query is empty",
			query: "",
			limit: 50,
			expected: []BusStop{
				{
					ID:          "01319",
					Description: "Lavender Stn Exit A/ICA",
					RoadName:    "Kallang Rd",
					Location:    Location{Latitude: 1.307574, Longitude: 103.86326},
				},
				{
					ID:          "01339",
					Description: "Bef Crawford Bridge",
					RoadName:    "Crawford St",
					Location:    Location{Latitude: 1.307746, Longitude: 103.864263},
				},
				{
					ID:          "07371",
					Description: "Aft Kallang Rd",
					RoadName:    "Lavender St",
					Location:    Location{Latitude: 1.309508, Longitude: 103.8635},
				},
			},
		},
		{
			name:  "returns up to limit bus stops when query is empty",
			query: "",
			limit: 2,
			expected: []BusStop{
				{
					ID:          "01319",
					Description: "Lavender Stn Exit A/ICA",
					RoadName:    "Kallang Rd",
					Location:    Location{Latitude: 1.307574, Longitude: 103.86326},
				},
				{
					ID:          "01339",
					Description: "Bef Crawford Bridge",
					RoadName:    "Crawford St",
					Location:    Location{Latitude: 1.307746, Longitude: 103.864263},
				},
			},
		},
		{
			name:  "doesn't treat slash as part of a token",
			query: "ica",
			limit: 5,
			expected: []BusStop{
				{
					ID:          "01319",
					Description: "Lavender Stn Exit A/ICA",
					RoadName:    "Kallang Rd",
					Location:    Location{Latitude: 1.307574, Longitude: 103.86326},
				},
			},
		},
	}
	repo := SQLBusStopRepository{DB: tx}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := repo.Search(tt.query, tt.limit)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
