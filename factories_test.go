package bete

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/yi-jiayu/datamall/v3"
)

var refTime = time.Unix(1584244383, 0)

func randomID() int {
	return rand.Int()
}

func randomInt64ID() int64 {
	return rand.Int63()
}

func randomStringID() string {
	p := make([]byte, 6)
	_, err := rand.Read(p)
	if err != nil {
		return "abcdef"
	}
	return base64.StdEncoding.EncodeToString(p)
}

func buildDataMallBusArrival() datamall.BusArrival {
	return datamall.BusArrival{
		BusStopCode: "96049",
		Services: []datamall.Service{
			{
				ServiceNo: "5",
				NextBus: datamall.ArrivingBus{
					EstimatedArrival: refTime.Add(-100 * time.Second),
					Load:             "SDA",
					Type:             "DD",
				},
				NextBus2: datamall.ArrivingBus{
					EstimatedArrival: refTime.Add(600 * time.Second),
					Load:             "SDA",
					Type:             "DD",
				},
				NextBus3: datamall.ArrivingBus{
					EstimatedArrival: refTime.Add(2200 * time.Second),
					Load:             "LSD",
					Feature:          "WAB",
					Type:             "BD",
				},
			},
			{
				ServiceNo: "24",
				NextBus: datamall.ArrivingBus{
					EstimatedArrival: refTime.Add(100 * time.Second),
					Load:             "SEA",
					Type:             "SD",
				},
				NextBus2: datamall.ArrivingBus{
					EstimatedArrival: refTime.Add(200 * time.Second),
					Load:             "SDA",
					Type:             "DD",
					Feature:          "WAB",
				},
				NextBus3: datamall.ArrivingBus{
					EstimatedArrival: refTime.Add(400 * time.Second),
					Load:             "LSD",
					Type:             "BD",
				},
			},
		},
	}
}

func buildBusStop() BusStop {
	return BusStop{
		ID:          "96049",
		Description: "UPP CHANGI STN/SUTD",
		RoadName:    "Upp Changi Rd East",
	}
}
