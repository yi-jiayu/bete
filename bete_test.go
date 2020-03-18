package bete

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate bin/mockgen -destination mocks_test.go -package bete -self_package github.com/yi-jiayu/bete . Clock,DataMall,Telegram,BusStopRepository,FavouriteRepository

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
	b, finish := newMockBete(t)
	defer finish()

	stop := buildBusStop()
	arrivals := buildDataMallBusArrival()

	b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
	b.BusStops.(*MockBusStopRepository).EXPECT().Find(gomock.Any()).Return(stop, nil)
	b.DataMall.(*MockDataMall).EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)

	actual, err := b.etaMessageText(context.Background(), stop.ID, nil)
	assert.NoError(t, err)
	expected, err := FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   nil,
	})
	if err != nil {
		panic(err)
	}
	assert.Equal(t, expected, actual)
}