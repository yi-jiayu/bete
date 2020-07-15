package bete

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yi-jiayu/datamall/v3"
)

func Test_sortByService(t *testing.T) {
	cases := []struct {
		name     string
		unsorted []datamall.Service
		sorted   []datamall.Service
	}{
		{
			name: "sorts services numerically",
			unsorted: []datamall.Service{
				{ServiceNo: "100"},
				{ServiceNo: "5"},
				{ServiceNo: "20"},
			},
			sorted: []datamall.Service{
				{ServiceNo: "5"},
				{ServiceNo: "20"},
				{ServiceNo: "100"},
			},
		},
		{
			name: "sorts services numerically and then by suffix",
			unsorted: []datamall.Service{
				{ServiceNo: "138B"},
				{ServiceNo: "138A"},
			},
			sorted: []datamall.Service{
				{ServiceNo: "138A"},
				{ServiceNo: "138B"},
			},
		},
		{
			name: "places services that do not start with a number at the end",
			unsorted: []datamall.Service{
				{ServiceNo: "NR1"},
				{ServiceNo: "20"},
			},
			sorted: []datamall.Service{
				{ServiceNo: "20"},
				{ServiceNo: "NR1"},
			},
		},
		{
			name: "sorts services that do not start with a number lexicographically",
			unsorted: []datamall.Service{
				{ServiceNo: "NR1"},
				{ServiceNo: "CT8"},
			},
			sorted: []datamall.Service{
				{ServiceNo: "CT8"},
				{ServiceNo: "NR1"},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got, want := sortByService(c.unsorted), c.sorted; !reflect.DeepEqual(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
		})
	}
}

func Test_filterByService(t *testing.T) {
	cases := []struct {
		name     string
		services []datamall.Service
		filter   []string
		expected []datamall.Service
	}{
		{
			name: "noop when filter is empty",
			services: []datamall.Service{
				{ServiceNo: "2"},
				{ServiceNo: "5"},
				{ServiceNo: "24"},
			},
			filter: nil,
			expected: []datamall.Service{
				{ServiceNo: "2"},
				{ServiceNo: "5"},
				{ServiceNo: "24"},
			},
		},
		{
			name: "returns only services in filter",
			services: []datamall.Service{
				{ServiceNo: "2"},
				{ServiceNo: "5"},
				{ServiceNo: "24"},
			},
			filter: []string{"5", "24"},
			expected: []datamall.Service{
				{ServiceNo: "5"},
				{ServiceNo: "24"},
			},
		},
		{
			name: "filter should be case-insensitive",
			services: []datamall.Service{
				{ServiceNo: "2A"},
				{ServiceNo: "5e"},
			},
			filter: []string{"2a", "5E"},
			expected: []datamall.Service{
				{ServiceNo: "2A"},
				{ServiceNo: "5e"},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			filtered := filterByService(c.filter, c.services)
			assert.Equal(t, c.expected, filtered)
		})
	}
}

func TestFormatArrivalsByService(t *testing.T) {
	cases := []struct {
		name     string
		arrivals ArrivalInfo
		expected string
	}{
		{
			name: "show bus stop details when available",
			arrivals: ArrivalInfo{
				Stop: BusStop{
					ID:          "96049",
					Description: "UPP CHANGI STN/SUTD",
					RoadName:    "Upp Changi Rd East",
				},
				Time:     refTime,
				Services: buildDataMallBusArrival().Services,
			},
			expected: `<strong>UPP CHANGI STN/SUTD (96049)</strong>
Upp Changi Rd East
<pre>
| Svc  | Nxt | 2nd | 3rd |
|------|-----|-----|-----|
| 5    |  -1 |  10 |  36 |
| 24   |   1 |   3 |   6 |
</pre>
<em>Last updated on Sun, 15 Mar 20 11:53 SGT</em>`,
		},
		{
			name: "escapes dangerous HTML",
			arrivals: ArrivalInfo{
				Stop: BusStop{
					ID: "<em>Hello</em>",
				},
				Time:     refTime,
				Services: buildDataMallBusArrival().Services,
				Filter:   []string{"2", "5", "<strong>24</strong>"},
			},
			expected: `<strong>&lt;em&gt;Hello&lt;/em&gt;</strong>
<pre>
| Svc  | Nxt | 2nd | 3rd |
|------|-----|-----|-----|
| 5    |  -1 |  10 |  36 |
</pre>
Filtered by services: 2, 5, &lt;strong&gt;24&lt;/strong&gt;
<em>Last updated on Sun, 15 Mar 20 11:53 SGT</em>`,
		},
		{
			name: "show only bus stop id when details not available",
			arrivals: ArrivalInfo{
				Stop: BusStop{
					ID: "96049",
				},
				Time:     refTime,
				Services: buildDataMallBusArrival().Services,
			},
			expected: `<strong>96049</strong>
<pre>
| Svc  | Nxt | 2nd | 3rd |
|------|-----|-----|-----|
| 5    |  -1 |  10 |  36 |
| 24   |   1 |   3 |   6 |
</pre>
<em>Last updated on Sun, 15 Mar 20 11:53 SGT</em>`,
		},
		{
			name: "filters services and shows filtered services when filter provided",
			arrivals: ArrivalInfo{
				Stop: BusStop{
					ID:          "96049",
					Description: "UPP CHANGI STN/SUTD",
					RoadName:    "Upp Changi Rd East",
				},
				Time:     refTime,
				Services: buildDataMallBusArrival().Services,
				Filter:   []string{"2", "24"},
			},
			expected: `<strong>UPP CHANGI STN/SUTD (96049)</strong>
Upp Changi Rd East
<pre>
| Svc  | Nxt | 2nd | 3rd |
|------|-----|-----|-----|
| 24   |   1 |   3 |   6 |
</pre>
Filtered by services: 2, 24
<em>Last updated on Sun, 15 Mar 20 11:53 SGT</em>`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := formatArrivalsSummary(c.arrivals)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestFormatArrivalsShowingDetails(t *testing.T) {
	cases := []struct {
		name     string
		arrivals ArrivalInfo
		expected string
	}{
		{
			name: "show bus stop details when available",
			arrivals: ArrivalInfo{
				Stop: BusStop{
					ID:          "96049",
					Description: "UPP CHANGI STN/SUTD",
					RoadName:    "Upp Changi Rd East",
				},
				Time:     refTime,
				Services: buildDataMallBusArrival().Services,
			},
			expected: `<strong>UPP CHANGI STN/SUTD (96049)</strong>
Upp Changi Rd East
<pre>
Svc   Eta  Sea  Typ  Fea
---   ---  ---  ---  ---
5      -1  SDA   DD     
24      1  SEA   SD     
24      3  SDA   DD  WAB
24      6  LSD   BD     
5      10  SDA   DD     
5      36  LSD   BD  WAB
</pre>
<em>Last updated on Sun, 15 Mar 20 11:53 SGT</em>`,
		},
		{
			name: "escapes dangerous HTML",
			arrivals: ArrivalInfo{
				Stop: BusStop{
					ID: "<em>Hello</em>",
				},
				Time:     refTime,
				Services: buildDataMallBusArrival().Services,
				Filter:   []string{"2", "5", "<strong>24</strong>"},
			},
			expected: `<strong>&lt;em&gt;Hello&lt;/em&gt;</strong>
<pre>
Svc   Eta  Sea  Typ  Fea
---   ---  ---  ---  ---
5      -1  SDA   DD     
5      10  SDA   DD     
5      36  LSD   BD  WAB
</pre>
Filtered by services: 2, 5, &lt;strong&gt;24&lt;/strong&gt;
<em>Last updated on Sun, 15 Mar 20 11:53 SGT</em>`,
		},
		{
			name: "show only bus stop id when details not available",
			arrivals: ArrivalInfo{
				Stop: BusStop{
					ID: "96049",
				},
				Time:     refTime,
				Services: buildDataMallBusArrival().Services,
			},
			expected: `<strong>96049</strong>
<pre>
Svc   Eta  Sea  Typ  Fea
---   ---  ---  ---  ---
5      -1  SDA   DD     
24      1  SEA   SD     
24      3  SDA   DD  WAB
24      6  LSD   BD     
5      10  SDA   DD     
5      36  LSD   BD  WAB
</pre>
<em>Last updated on Sun, 15 Mar 20 11:53 SGT</em>`,
		},
		{
			name: "filters services and shows filtered services when filter provided",
			arrivals: ArrivalInfo{
				Stop: BusStop{
					ID:          "96049",
					Description: "UPP CHANGI STN/SUTD",
					RoadName:    "Upp Changi Rd East",
				},
				Time:     refTime,
				Services: buildDataMallBusArrival().Services,
				Filter:   []string{"2", "24"},
			},
			expected: `<strong>UPP CHANGI STN/SUTD (96049)</strong>
Upp Changi Rd East
<pre>
Svc   Eta  Sea  Typ  Fea
---   ---  ---  ---  ---
24      1  SEA   SD     
24      3  SDA   DD  WAB
24      6  LSD   BD     
</pre>
Filtered by services: 2, 24
<em>Last updated on Sun, 15 Mar 20 11:53 SGT</em>`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := formatArrivalsDetails(c.arrivals)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestFormatArrivals(t *testing.T) {
	info := ArrivalInfo{
		Stop: BusStop{
			ID:          "96049",
			Description: "UPP CHANGI STN/SUTD",
			RoadName:    "Upp Changi Rd East",
		},
		Time:     refTime,
		Services: buildDataMallBusArrival().Services,
	}
	tests := []struct {
		name     string
		format   Format
		expected string
	}{
		{
			name:     "summary",
			format:   FormatSummary,
			expected: must(formatArrivalsSummary(info)).(string),
		},
		{
			name:     "details",
			format:   FormatDetails,
			expected: must(formatArrivalsDetails(info)).(string),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := FormatArrivals(info, tt.format)
			assert.NoError(t, err)
			assert.Equal(t, actual, tt.expected)
		})
	}
}
