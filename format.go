package bete

import (
	"bytes"
	"html/template"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yi-jiayu/datamall/v3"
)

var trimTrailingLettersRegexp = regexp.MustCompile("[^0-9]$")

var funcMap = map[string]interface{}{
	"arrivingBuses":   arrivingBuses,
	"filterByService": filterByService,
	"inSGT":           inSGT,
	"join":            strings.Join,
	"sortByArrival":   sortByArrival,
	"sortByService":   sortByService,
	"take":            take,
	"until":           minutesUntil,
}

var (
	arrivalSummaryTemplate = template.Must(template.New("arrival_summary").
				Funcs(funcMap).
				Parse(templateArrivalSummary))
	arrivalDetailsTemplate = template.Must(template.New("arrival_details").
				Funcs(funcMap).
				Parse(templateArrivalDetails))
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
			if strings.EqualFold(s.ServiceNo, f) {
				filtered = append(filtered, s)
			}
		}
	}
	return filtered
}

type ArrivingBus struct {
	ServiceNo string
	datamall.ArrivingBus
}

func arrivingBuses(services []datamall.Service) []ArrivingBus {
	var buses []ArrivingBus
	for _, service := range services {
		if !service.NextBus.EstimatedArrival.IsZero() {
			buses = append(buses, ArrivingBus{
				ServiceNo:   service.ServiceNo,
				ArrivingBus: service.NextBus,
			})
		}
		if !service.NextBus2.EstimatedArrival.IsZero() {
			buses = append(buses, ArrivingBus{
				ServiceNo:   service.ServiceNo,
				ArrivingBus: service.NextBus2,
			})
		}
		if !service.NextBus3.EstimatedArrival.IsZero() {
			buses = append(buses, ArrivingBus{
				ServiceNo:   service.ServiceNo,
				ArrivingBus: service.NextBus3,
			})
		}
	}
	return buses
}

func sortByArrival(buses []ArrivingBus) []ArrivingBus {
	sort.Slice(buses, func(i, j int) bool {
		return buses[i].EstimatedArrival.Before(buses[j].EstimatedArrival)
	})
	return buses
}

// take returns up to the first n arriving buses.
func take(n int, arriving []ArrivingBus) []ArrivingBus {
	if n <= 0 || n > len(arriving) {
		n = len(arriving)
	}
	return arriving[:n]
}

type ArrivalInfo struct {
	Stop     BusStop
	Time     time.Time
	Services []datamall.Service
	Filter   []string
}

func FormatArrivalsByService(arrivals ArrivalInfo) (string, error) {
	b := new(bytes.Buffer)
	err := arrivalSummaryTemplate.Execute(b, arrivals)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func FormatArrivalsShowingDetails(arrivals ArrivalInfo) (string, error) {
	b := new(bytes.Buffer)
	err := arrivalDetailsTemplate.Execute(b, arrivals)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
