package bete

import (
	"bytes"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/yi-jiayu/datamall/v3"
)

var trimTrailingLettersRegexp = regexp.MustCompile("[^0-9]$")

var funcMap = map[string]interface{}{
	"inSGT":           inSGT,
	"filterByService": filterByService,
	"join":            strings.Join,
	"sortByService":   sortByService,
	"until":           minutesUntil,
}

var (
	arrivalsByServiceTemplate = template.Must(template.New("arrivals_by_service.tmpl").
		Funcs(funcMap).
		ParseFiles("templates/arrivals_by_service.tmpl"))
)

var sgt = time.FixedZone("SGT", 8*3600)

func inSGT(t time.Time) string {
	return t.In(sgt).Format("Mon, 02 Jan 06 15:04 MST")
}

// minutesUntil returns the number of minutes from now until then.
func minutesUntil(now time.Time, then time.Time) string {
	if then.IsZero() {
		return "?"
	}
	return strconv.Itoa(int(then.Sub(now).Minutes()))
}

func sortByService(services []datamall.Service) []datamall.Service {
	sort.Slice(services, func(i, j int) bool {
		first, err1 := strconv.Atoi(trimTrailingLettersRegexp.ReplaceAllString(services[i].ServiceNo, ""))
		second, err2 := strconv.Atoi(trimTrailingLettersRegexp.ReplaceAllString(services[j].ServiceNo, ""))
		switch {
		case err1 != nil && err2 != nil:
			// when both services cannot be parsed as integers, sort them lexicographically
			return services[i].ServiceNo < services[j].ServiceNo
		case err1 == nil && err2 != nil:
			// if j cannot be parsed as an integer, then i should come before j
			return true
		case err1 != nil && err2 == nil:
			// if i cannot be parsed as an integer, then j should come before i
			return false
		}
		if first == second {
			return services[i].ServiceNo < services[j].ServiceNo
		}
		return first < second
	})
	return services
}

func filterByService(filter []string, services []datamall.Service) []datamall.Service {
	if len(filter) == 0 {
		return services
	}
	var filtered []datamall.Service
	for _, s := range services {
		for _, f := range filter {
			if s.ServiceNo == f {
				filtered = append(filtered, s)
			}
		}
	}
	return filtered
}

type ArrivalInfo struct {
	Stop     BusStop
	Time     time.Time
	Services []datamall.Service
	Filter   []string
}

func FormatArrivalsByService(arrivals ArrivalInfo) (string, error) {
	b := new(bytes.Buffer)
	err := arrivalsByServiceTemplate.Execute(b, arrivals)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
