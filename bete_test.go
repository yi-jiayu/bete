package bete

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yi-jiayu/datamall/v3"
)

//go:generate bin/mockgen -destination mocks_test.go -package bete -self_package github.com/yi-jiayu/bete . Clock,DataMall,Telegram,BusStopRepository,FavouriteRepository

func init() {
	// Disable logging in tests.
	log.SetOutput(ioutil.Discard)
}

func must(i interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return i
}

func newMockBete(t *testing.T) (Bete, func()) {
	ctrl := gomock.NewController(t)
	b := Bete{
		Clock:      NewMockClock(ctrl),
		BusStops:   NewMockBusStopRepository(ctrl),
		Favourites: NewMockFavouriteRepository(ctrl),
		DataMall:   NewMockDataMall(ctrl),
		Telegram:   NewMockTelegram(ctrl),
	}
	return b, ctrl.Finish
}

func TestBete_etaMessageText(t *testing.T) {
	stop := buildBusStop()
	arrivals := buildDataMallBusArrival()
	tests := []struct {
		name     string
		format   Format
		arrivals datamall.BusArrival
		err      error
		errMsg   string
	}{
		{
			name:     "summary",
			format:   FormatSummary,
			arrivals: arrivals,
		},
		{
			name:     "details",
			format:   FormatDetails,
			arrivals: arrivals,
		},
		{
			name:   "datamall error",
			err:    &datamall.Error{StatusCode: http.StatusNotFound},
			errMsg: `Error getting bus arrivals from LTA DataMall (HTTP status 404)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, finish := newMockBete(t)
			defer finish()

			b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
			b.BusStops.(*MockBusStopRepository).EXPECT().Find(stop.ID).Return(stop, nil)
			b.DataMall.(*MockDataMall).EXPECT().GetBusArrival(stop.ID, "").Return(tt.arrivals, tt.err)

			actual, err := b.etaMessageText(context.Background(), stop.ID, nil, tt.format)
			assert.NoError(t, err)
			expected := must(FormatArrivals(ArrivalInfo{
				Stop:     stop,
				Time:     refTime,
				Services: tt.arrivals.Services,
				Filter:   nil,
				ErrMsg:   tt.errMsg,
			}, tt.format)).(string)
			assert.Equal(t, expected, actual)
		})
	}
}
