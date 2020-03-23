package bete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getStreetViewStaticURL(t *testing.T) {
	key := "API_KEY"
	stop := buildBusStop()

	actual := getStreetViewStaticURL(key, stop)
	expected := "https://maps.googleapis.com/maps/api/streetview?key=API_KEY&location=1.340874%2C103.961433&size=100x100"
	assert.Equal(t, expected, actual)
}
