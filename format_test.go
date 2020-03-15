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

func TestFormatArrivalsByService(t *testing.T) {
	t.Run("without filter", func(t *testing.T) {
		arrivals := ArrivalInfo{
			Stop: BusStop{
				ID:          "96049",
				Description: "UPP CHANGI STN/SUTD",
				RoadName:    "Upp Changi Rd East",
			},
			Time:     refTime,
			Services: buildDataMallBusArrival().Services,
		}
		formatted, err := FormatArrivalsByService(arrivals)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, `<strong>UPP CHANGI STN/SUTD (96049)</strong>
Upp Changi Rd East
<pre>
| Svc  | Nxt | 2nd | 3rd |
|------|-----|-----|-----|
| 5    |  -1 |  10 |  36 |
| 24   |   1 |   3 |   6 |
</pre>
<em>Last updated on Sun, 15 Mar 20 11:53 SGT</em>`, formatted)
	})
	t.Run("with filter", func(t *testing.T) {
		arrivals := ArrivalInfo{
			Stop: BusStop{
				ID:          "96049",
				Description: "UPP CHANGI STN/SUTD",
				RoadName:    "Upp Changi Rd East",
			},
			Time:     refTime,
			Services: buildDataMallBusArrival().Services,
			Filter:   []string{"2", "24"},
		}
		formatted, err := FormatArrivalsByService(arrivals)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, `<strong>UPP CHANGI STN/SUTD (96049)</strong>
Upp Changi Rd East
<pre>
| Svc  | Nxt | 2nd | 3rd |
|------|-----|-----|-----|
| 24   |   1 |   3 |   6 |
</pre>
Filtered by services: 2, 24
<em>Last updated on Sun, 15 Mar 20 11:53 SGT</em>`, formatted)
	})
}

func Test_filterByService(t *testing.T) {
	t.Run("noop when filter is empty", func(t *testing.T) {
		services := []datamall.Service{
			{ServiceNo: "2"},
			{ServiceNo: "5"},
			{ServiceNo: "24"},
		}
		filtered := filterByService(nil, services)
		assert.Equal(t, []datamall.Service{
			{ServiceNo: "2"},
			{ServiceNo: "5"},
			{ServiceNo: "24"},
		}, filtered)
	})
	t.Run("returns only services in filter", func(t *testing.T) {
		services := []datamall.Service{
			{ServiceNo: "2"},
			{ServiceNo: "5"},
			{ServiceNo: "24"},
		}
		filter := []string{"5", "24"}
		filtered := filterByService(filter, services)
		assert.Equal(t, []datamall.Service{
			{ServiceNo: "5"},
			{ServiceNo: "24"},
		}, filtered)
	})
}
