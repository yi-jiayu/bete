package bete

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/yi-jiayu/ted"
)

func TestBete_HandleMessage_Text(t *testing.T) {
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
	filter := []string{"5", "24"}
	arrivals := buildDataMallBusArrival()
	chatID := randomID()
	text := must(FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   filter,
	})).(string)
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stop.ID, filter),
	}

	clock.EXPECT().Now().Return(refTime)
	busStopRepository.EXPECT().Find(gomock.Any()).Return(stop, nil)
	dm.EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	telegram.EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			Chat: ted.Chat{ID: chatID},
			Text: "96049 5 24",
		},
	}
	b.HandleUpdate(update)
}
