package bete

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yi-jiayu/ted"
)

func TestBete_HandleCallbackQuery_Refresh(t *testing.T) {
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

	clock.EXPECT().Now().Return(refTime)
	busStopRepository.EXPECT().Find(gomock.Any()).Return(stop, nil)
	dm.EXPECT().GetBusArrival(stop.ID, "").Return(arrivals, nil)
	telegram.EXPECT().Do(editMessageText).Return(ted.Response{}, nil)
	telegram.EXPECT().Do(answerCallbackQuery).Return(ted.Response{}, nil)

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
	err := b.HandleUpdate(update)
	assert.NoError(t, err)
}
