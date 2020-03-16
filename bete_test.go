package bete

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yi-jiayu/ted"
)

//go:generate bin/mockgen -source bete.go -destination bete_mocks_test.go -package bete

func must(i interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return i
}

func TestBete_etaMessageText(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clock := NewMockClock(ctrl)
	busStopRepository := NewMockBusStopRepository(ctrl)
	dm := NewMockDataMall(ctrl)
	b := Bete{
		Clock:    clock,
		BusStops: busStopRepository,
		DataMall: dm,
	}

	stop := buildBusStop()
	arrivals := buildDataMallBusArrival()

	clock.EXPECT().Now().Return(refTime)
	busStopRepository.EXPECT().Find(gomock.Any()).Return(stop, nil)
	dm.EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)

	actual, err := b.etaMessageText(stop.ID, nil)
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

func TestBete_SendETAMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clock := NewMockClock(ctrl)
	busStopRepository := NewMockBusStopRepository(ctrl)
	dm := NewMockDataMall(ctrl)
	telegram := NewMockTelegram(ctrl)
	b := Bete{
		Clock:    clock,
		BusStops: busStopRepository,
		DataMall: dm,
		Telegram: telegram,
	}

	stop := buildBusStop()
	arrivals := buildDataMallBusArrival()
	chatID := 123
	text := must(FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   nil,
	})).(string)
	req := ted.SendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "HTML",
	}

	clock.EXPECT().Now().Return(refTime)
	busStopRepository.EXPECT().Find(gomock.Any()).Return(stop, nil)
	dm.EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	telegram.EXPECT().Do(req).Return(ted.Response{}, nil)

	err := b.SendETAMessage(chatID, stop.ID, nil)
	assert.NoError(t, err)
}
