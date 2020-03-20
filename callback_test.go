package bete

import (
	"context"
	"fmt"
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
	chatID := randomInt64ID()
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
				Type:   callbackRefresh,
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
	chatID := randomInt64ID()
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
				Type:   callbackResend,
				StopID: stop.ID,
				Filter: filter,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_HandleCallbackQuery_AddFavourite(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	sendMessage := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        stringAddFavouritePromptForQuery,
		ReplyMarkup: ted.ForceReply{},
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
	}

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
				Type: callbackAddFavourite,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_saveFavouriteCallback_WithName(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	name := "Opp Tropicana Condo"
	query := Query{Stop: "96049", Filter: []string{"5", "24"}}
	favourites := []string{"Home", name}
	showFavourites := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        fmt.Sprintf("Added the query %q to your favourites as %q!", query.Canonical(), name),
		ReplyMarkup: showFavouritesReplyMarkup(favourites),
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
	}
	removeButtons := ted.EditMessageReplyMarkupRequest{
		ChatID:    chatID,
		MessageID: messageID,
	}

	b.Favourites.(*MockFavouriteRepository).EXPECT().Put(userID, name, query.Canonical())
	b.Favourites.(*MockFavouriteRepository).EXPECT().List(userID).Return(favourites, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(showFavourites).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallbackQuery).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(removeButtons).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID:   callbackQueryID,
			From: ted.User{ID: userID},
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type:  callbackSaveFavourite,
				Query: &query,
				Name:  name,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}

func TestBete_saveFavouriteCallback_WithoutName(t *testing.T) {
	b, finish := newMockBete(t)
	defer finish()

	userID := randomID()
	chatID := randomInt64ID()
	messageID := randomID()
	callbackQueryID := randomStringID()
	query := Query{Stop: "96049", Filter: []string{"5", "24"}}
	promptForName := ted.SendMessageRequest{
		ChatID:      chatID,
		Text:        fmt.Sprintf(stringAddFavouritePromptForName, query.Canonical()),
		ReplyMarkup: ted.ForceReply{},
	}
	answerCallbackQuery := ted.AnswerCallbackQueryRequest{
		CallbackQueryID: callbackQueryID,
	}
	removeButtons := ted.EditMessageReplyMarkupRequest{
		ChatID:    chatID,
		MessageID: messageID,
	}

	b.Telegram.(*MockTelegram).EXPECT().Do(promptForName).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(answerCallbackQuery).Return(ted.Response{}, nil)
	b.Telegram.(*MockTelegram).EXPECT().Do(removeButtons).Return(ted.Response{}, nil)

	update := ted.Update{
		CallbackQuery: &ted.CallbackQuery{
			ID:   callbackQueryID,
			From: ted.User{ID: userID},
			Message: &ted.Message{
				ID:   messageID,
				Chat: ted.Chat{ID: chatID},
			},
			Data: CallbackData{
				Type:  callbackSaveFavourite,
				Query: &query,
			}.Encode(),
		},
	}
	b.HandleUpdate(context.Background(), update)
}
