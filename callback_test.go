package bete

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/yi-jiayu/ted"
)

func TestBete_HandleCallbackQuery_Refresh(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	stop := buildBusStop()
	filter := []string{"5", "24"}
	arrivals := buildDataMallBusArrival()
	chatID := randomID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	text := must(FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   filter,
	})).(string)
	editMessageText := ted.EditMessageTextRequest{
		ChatID:      chatID,
		MessageID:   messageID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkupP(stop.ID, filter),
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            "ETAs updated!",
	}

	b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
	b.BusStops.(*MockBusStopRepository).EXPECT().Find(gomock.Any()).Return(stop, nil)
	b.DataMall.(*MockDataMall).EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(editMessageText).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallbackQuery).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID: callbackQueryID,
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type:   "refresh",
				StopID: stop.ID,
				Filter: filter,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleCallbackQuery_Resend(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	stop := buildBusStop()
	filter := []string{"5", "24"}
	arrivals := buildDataMallBusArrival()
	chatID := randomID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	text := must(FormatArrivalsByService(ArrivalInfo{
		Stop:     stop,
		Time:     refTime,
		Services: arrivals.Services,
		Filter:   filter,
	})).(string)
	sendMessage := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        text,
		ParseMode:   "HTML",
		ReplyMarkup: etaMessageReplyMarkup(stop.ID, filter),
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
		Text:            "ETAs sent!",
	}

	b.Clock.(*MockClock).EXPECT().Now().Return(refTime)
	b.BusStops.(*MockBusStopRepository).EXPECT().Find(gomock.Any()).Return(stop, nil)
	b.DataMall.(*MockDataMall).EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(sendMessage).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallbackQuery).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID: callbackQueryID,
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type:   "resend",
				StopID: stop.ID,
				Filter: filter,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}
