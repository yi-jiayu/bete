package bete

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/yi-jiayu/ted"
)

func TestBete_HandleTextMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clock := NewMockClock(ctrl)
	busStopRepository := NewMockBusStopRepository(ctrl)
	favouriteRepository := NewMockFavouriteRepository(ctrl)
	dm := NewMockDataMall(ctrl)
	telegram := NewMockTelegram(ctrl)
	b := Bete{
		Clock:      clock,
		BusStops:   busStopRepository,
		Favourites: favouriteRepository,
		DataMall:   dm,
		Telegram:   telegram,
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
	favouriteRepository.EXPECT().FindByUserAndText(gomock.Any(), gomock.Any()).Return("")
	dm.EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	telegram.EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: randomID()},
			Chat: ted.Chat{ID: chatID},
			Text: "96049 5 24",
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleTextMessage_Favourite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clock := NewMockClock(ctrl)
	busStopRepository := NewMockBusStopRepository(ctrl)
	favouriteRepository := NewMockFavouriteRepository(ctrl)
	dm := NewMockDataMall(ctrl)
	telegram := NewMockTelegram(ctrl)
	b := Bete{
		Clock:      clock,
		BusStops:   busStopRepository,
		Favourites: favouriteRepository,
		DataMall:   dm,
		Telegram:   telegram,
	}

	stop := buildBusStop()
	filter := []string{"5", "24"}
	arrivals := buildDataMallBusArrival()
	userID := randomID()
	chatID := randomID()
	messageText := "SUTD"
	replyText := must(FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   filter,
	})).(string)
	req := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        replyText,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stop.ID, filter),
	}

	clock.EXPECT().Now().Return(refTime)
	busStopRepository.EXPECT().Find(gomock.Any()).Return(stop, nil)
	favouriteRepository.EXPECT().FindByUserAndText(userID, messageText).Return("96049 5 24")
	dm.EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	telegram.EXPECT().Do(req).Return(ted.Response{}, nil)

	update := ted.Update{
		Message: &ted.Message{
			From: &ted.User{ID: userID},
			Chat: ted.Chat{ID: chatID},
			Text: messageText,
		},
	}
	b.HandleUpdate(context.Background(), update)
}
