package bete

import (
	"fmt"
	"net/url"
)

// endpointStreetViewStaticAPI is the endpoint for the Street View Image API.
const endpointStreetViewStaticAPI = "https://maps.googleapis.com/maps/api/streetview"

func getStreetViewStaticURL(key string, stop BusStop) string {
	params := url.Values{}
	params.Set("key", key)
	params.Set("location", fmt.Sprintf("%f,%f", stop.Location.Latitude, stop.Location.Longitude))
	params.Set("size", fmt.Sprintf("%dx%d", 100, 100))

	return fmt.Sprintf("%s?%s", endpointStreetViewStaticAPI, params.Encode())
}
