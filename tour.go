package bete

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Tour section slugs
const (
	tourSectionStart = "start"
)

// TourSectionData contains what to display for a tour section.
type TourSectionData struct {
	Text       string                  `yaml:"text"`
	Navigation []TourSectionNavigation `yaml:"navigation"`
}

// TourSectionNavigation describes a link between tour sections.
type TourSectionNavigation struct {
	Text   string `yaml:"text"`
	Target string `yaml:"target"`
}

var tour = mustLoadTour("tour.yaml")

func mustLoadTour(path string) map[string]TourSectionData {
	var tour map[string]TourSectionData
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	err = yaml.NewDecoder(f).Decode(&tour)
	if err != nil {
		panic(err)
	}
	return tour
}
